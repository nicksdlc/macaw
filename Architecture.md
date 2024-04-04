# Initialization
The main entity is Context.
Context has Communicatior and runner.
There are 2 runner: Receiver and Sender.

# Receiver mode
In Sender mode responder.NewMessageResponder is created. It has references to Responses andcommunicator
Before running Responder builds ResponsePrototypes by calling prototype.NewResponsePrototypeBuilder. 
`ResponsePrototypeBuildeer` transforms all responses from Config into MessagePrototypes. Message prototypes have alias, templates for body and matchers. 
HTTP Comminicator in method RespondWith creates handlers for each MessagePrototype. These handlers are stored in map by URL. Handler have same signature as http.HandlerFunc, so can be used in satandard Mux.



# Communicator
Communicatior is used in Sender and Receiver modes. There are 2 types of communicators currently: HTTP and RMQ.


## Http Communicator

