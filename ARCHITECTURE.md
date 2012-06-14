# ARCHITECTURE

Uses Redis as a datastore, Go as a basic language. It is meant to be
heavily multithreading, with a moderately decentralized architecture,
allowing it to be run as a cluster.

## API

API client process. RESTful methods to administer the server as well
as to pull data for the UI.

## Control

Singleton thread, which handles the following tasks:

 * Database schema control and patching
 * Scheduling checks

Should attempt to "lock" based on the first server which contacts
the database backend.

## Poll

Threads which use BRPOPLPUSH to wait for scheduled checks, then
execute those checks, and perform some logic to determine whether or
not they need alerting.

## Alert

Alerting and escalation, using BRPOPLPUSH to wait for notifications,
then follow a particular escalation pattern based on the notifications
which need to be sent.

## Data Structures

### Checks

swarm:index:checks SET - all checks
swarm:checks:KEY
swarm:checks:KEY:hosts SET- consists of host keys
swarm:checks:KEY:groups SET - consist of group keys

### Hosts

swarm:hosts:KEY

### Groups

swarm:groups:KEY

### Contacts

swarm:contacts:KEY

## Queues

swarm:queue:alert SET - active alert/notification queue
swarm:queue:poll SET - active poll queue

