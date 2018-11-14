# Byzantine Generals Problem

## Usage

`go run main.go`

This will randomly choose faulty generals and display the consensus decision at the end.

## Design

In distributed systems, there is often a problem of determining what the general agreement is among Nodes of the system. It is possible that some nodes may fail or be faulty but the group of nodes that are not faulty should agree upon a decision. The idea of the algorithm is that depending on how many faulty nodes we have within our distributed system, we will need to send messages that many times to ensure that the faulty messages do not cause a faulty decision. 

This algorithm guarantees that the generals will come to an agreement on the correct decision. There is one exception to this algorithm and that is when the commander(or first general) is one of the faulty generals. In this case, the node will send faulty messages to half of the nodes because it chooses faulty messages based off the id of the node it is sending too. In addition, this algorithm only guarantees that it will produce an agreed decision if the number of faulty generals is less than a third of the total number of generals. One major variation between this algorithm and some of the other ones online is the decision to make faulty generals not 100% faulty. In this algorithm a faulty general will only send a faulty message if the general it is sending it to has an even ID. This can change the output of the algorithm because it means there is overall less faulty messages. 