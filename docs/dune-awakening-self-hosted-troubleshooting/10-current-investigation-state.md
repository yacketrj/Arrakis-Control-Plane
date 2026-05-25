# Current Investigation State

This document preserves the current working investigation state for the active Dune: Awakening self-hosted troubleshooting effort.

## Evidence Boundary

Use only evidence from the uploaded `logs.zip`, files uploaded after that point, and user-provided prompts/messages after that upload. Do not import facts, assumptions, client details, player details, or conclusions from unrelated labels or older projects.

## Known Environment From User-Provided Evidence

```text
Hosting stack: Proxmox bare-metal host, Ubuntu 24.04 VM, AMP control panel, Docker behind AMP
CPU allocation: 16 CPU cores from AMD Ryzen 9 9950X3D
Memory allocation: 96 GB assigned to VM from 192 GB host memory
Network: data center hosting, 5 Gbps fiber, 5 static IPs
Example AMP instance path: /home/amp/.ampdata/instances/DuneAwakening01/
```

## Prior Working Findings

```text
The process-visible OS shown in game logs may reflect container or guest runtime context and should not be treated as proof of the physical host OS.
The environment owner confirmed the broader host stack separately: Proxmox bare metal, Ubuntu VM, AMP, Docker behind AMP.
A prior UID/GID mismatch caused UserEngine.ini write failures when host path ownership did not match the container's game user context.
RabbitMQ logs previously showed per-player queue delete failures where the queue was still in use.
Prior director logs showed many ClassicalInstancing groups with empty server lists, while known primary partitions such as Survival_1 and Overmap reported server IDs.
```

## Latest Uploaded Evidence: 2026-05-25 Travel Attempt

A newly provided log excerpt shows a failed travel attempt from Overmap/Overland to `Story_ArtOfKanly` / partition 13.

Sensitive values from the uploaded log, including server password and player/account identifiers, are intentionally not reproduced here.

### Observed Flow

```text
06:49:00 - Player is reported on Overland / Overmap context.
06:49:00 - Travel initiated to Story_ArtOfKanly_DuelingRing_Lobby.
06:49:00 - Travel Initialize stage reports success.
06:49:01 - Travel Request stage reports success.
06:49:01 - TravelQueue Story_ArtOfKanly receives the request with DestinationPartitionId=13 and OriginId=Overmap2.
06:49:01 - Director receives one-player travel request for Story_ArtOfKanly with instancingMode=ClassicalInstancing.
06:49:01 - Director creates a travel request and assigns queue token 5.
06:49:01 - Director returns TravelResponse Code=Queued, ServerState=Ready, ServerFull=True, DestinationPartitionId=null, ServerLoginToken empty, Queue length 1.
06:49:02 - Story_ArtOfKanly / partition 13 reports ready=true with no players.
06:49:04 - Client/source-side travel fails with Bgd Travel Failed! Error:5.
06:49:04 - Same request ID appears again in TravelQueue Story_ArtOfKanly, but OriginId changes to Nowhere.
```

## Current Interpretation

The latest evidence strengthens the theory that the failure is occurring inside the director/instancing travel allocation path before a usable destination login handoff is issued to the player.

The most important anomaly is that the destination partition reports ready, but the travel response still returns `ServerFull=True`, `DestinationPartitionId=null`, and an empty `ServerLoginToken`. That means the source/client appears not to receive a concrete usable destination even though the destination's server state exists.

The second important anomaly is the repeated/same request ID appearing with `OriginId=Nowhere` after the source-side travel failure. This may indicate orphaned, retried, or cleanup-path queue handling after the original Overmap-originated request fails.

## Current Leading Hypotheses

```text
1. ClassicalInstancing allocation/capacity state mismatch:
   The director sees the destination server state as ready but still treats the destination as full or unavailable for handoff.

2. Travel queue handoff failure:
   The request reaches the correct queue and receives a queue token, but the director does not issue a destination partition and login token before the source travel request fails.

3. Stale or inconsistent instancing metadata:
   Prior evidence showed ClassicalInstancing groups with empty server lists. The latest evidence shows partition 13 reporting ready, but not being allocated as a usable destination.

4. Post-failure orphan/retry path:
   The same request ID reappearing with OriginId=Nowhere after Error:5 suggests cleanup or retry logic may be re-queuing or re-processing the failed request without its original source context.
```

## Evidence Needed Next

```text
1. Destination map logs for the Story_ArtOfKanly container/service covering 06:48:45-06:49:15 UTC.
2. Director logs filtered by the request ID around 06:48:45-06:49:15 UTC.
3. RabbitMQ queue/consumer state captured during a controlled retry.
4. Active UDP listener output before/during/after the retry.
5. Current director/partition metadata showing map, partition_id, server_id, dimension_index, blocked, and label.
6. Any director DB or config values related to ClassicalInstancing capacity, hard cap, queue size, and partition assignment for Story_ArtOfKanly.
```

## Next Recommended Controlled Test

Run one controlled travel attempt to `Story_ArtOfKanly` while capturing:

```text
Director logs
Overmap/source logs
Story_ArtOfKanly destination logs
RabbitMQ queue/consumer snapshot
UDP listener output
Packet capture if possible
```

The target is to determine whether the destination ever receives a player login/prelogin/authorization attempt. If it does not, the failure is likely before destination connection and remains in director allocation, queue response, token issuance, or source handoff.
