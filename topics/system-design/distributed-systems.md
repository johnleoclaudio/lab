# Distributed Systems

### Finished Lectures:
- [x] Distributed Systems 1.1: Introduction

- Throwing networking into the mix together with your computing
- Multiple nodes communicating via a network trying to achieve some task together
- How we coordinate these nodes to achieve a goal

### What makes a system distributed?
- It's inherently distributed
  - sms for example
- For better reliability
  - horizontally scaled nodes helps availability and reliability
- For better performance:
  - get data from nearby node instead of getting data from US while in the philippines
- To solve bigger problem that then can with single computer 

### Why NOT make a system distributed?
- Network communication is unreliable
- wifi weak, cellsite far
- communication may fail 
- processes may crash and we even may not know!
- all this happen nondeterministically

### Fault Tolerance - The Goal
- Hard to do 
- we want the system to tolerate fault
- "if you can solve using a single computer, that's easy"
