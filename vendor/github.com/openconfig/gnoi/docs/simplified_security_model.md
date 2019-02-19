# Simplified security model for the gNOI Certificate service

Contributors: Sam Ribeiro (ribeiro@google.com), Rob Shakir (robjs@google.com),
Alireza Ghaffarkhah (aghaffar@google.com), Eric Breverman (ejbrever@google.com),
Anees Shaikh (aashaikh@google.com)

Updated: August 7th, 2018

Version: 0.1.0

## Simplified security model

What follows is the description of one security model and its relation with the
gNOI Certificate Service. This security model does not assert itself as the only
one that can be employed. Different security models can be defined to address
specific scenarios.

This simplified security model makes assumptions about the installation and
rotation of Certificates & CA Certificates. Their assumed usage is applicable
not only to gNOI and gNMI, but also other services like OpenFlow.

This simplified security model presumes that a secure connection between the
Target and the Client already exist for installing and rotating Certificates.
Bootstrapping is out of scope for this model.

This simplified security model assumes low risk of man-in-the-middle attacks
because it relinquishes mandating DNS verification for the Common Name of the
respective Client and Target Certificates.

## Client

A Client is defined as per the definition in the gNMI and gNOI specifications.

## Target

A Target is defined as per the definition in the gNMI and gNOI specifications.

## Peer

Peer is an entity that participates in establishing a connection with another
entity, using whatever protocol that makes use of Certificates to secure that
connection. An example would be an OpenFlow connection between an OpenFlow
switch and an OpenFlow Controller. Target and Client are particular cases of a
Peer.

## Client is a Relying Party

It is assumed that the Certificate Authority (CA) provides certificate signing
and generation services to the Client. The Client is therefore a Relying Party
for both certificates and CA bundles. The relationship and interaction between
the CA and the Client is out of scope for this model.

## Secure Connection

A secure connection is established between two Peers when mutual authentication
exists. Meaning that both Peers identify each other as genuine. An example is a
gNOI connection between the Client and the Target. The procedure for mutual
authentication is described below.

## Target is authorized trust by the Client

By authorized trust it is meant that the Target is granted the necessary
credentials (by the Client) to be able to identify Peers as genuine (using
whatever service uses these Certificates). The Target is authorized trust by the
Client via the successful execution of either:

*   gNOI certificate installation or rotation under a Secure connection (where
    the Target is not compromised), or
*   other provisioning or bootstrapping mechanism described outside of the scope
    of this security model.

## Target Certificate installation or rotation

To initiate certificate installation using the gNOI Certificate Management
service:

1.  The Client sends parameters to the target so that it can generate a
    Certificate Signing Request (CSR).
2.  The Target then generates this CSR using its Private Key.
3.  The Target then replies with the Certificate to the Client, which
    subsequently provides the Certificate to the CA to be signed.
4.  The signed Certificate is then returned to the Target by the Client and it
    becomes one of the Target’s Certificates.

Each one of the Target’s Certificate has a unique identifier, the Certificate
ID. This ID is used both to install and to rotate a Certificate. Certificate
installation for an existing Certificate ID MUST fail. Replacing an existing
Certificate must rather use the Certificate rotation mechanism, whose steps are
similar to the ones described above with the addition of final validation.

## Target CA pool

In order for the Target to validate the Client or a Peer's Certificates, it must
have a pool of one or more CA Certificates. These are provisioned onto the
Target by the Client during creation or rotation of the Target’s Certificates.

## Mutual Authentication

Mutual authentication exists when two Peers validate each other's Certificate
against their CA pools. Optionally, either the Client or the Target can validate
that the Common Name (CN) in each other’s Certificates matches the resolved one
for the Peer address of the connection.
