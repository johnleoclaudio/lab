Assignment
While Pub/Sub is an architecture that often refers to messages being sent between programs, or sometimes even between machines on a network, let's start with a simple pub/sub system that runs in a single program.

In "Peril", users control armies, represented as an array of pieces. When a piece moves to a location that contains a piece owned by another player, a battle occurs. Whenever a player makes a "move", we need to publish a message telling the game that a move has occurred.

Read the doBattles function. Note that it expects a list of Move messages that were published earlier.

Complete the march function. It should create a move message and send it via the publishMove function.
