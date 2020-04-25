# gsql
GSQL - Graph SQL, An alteration and facade of the SQL language for querying Graph Models.


## Overview
We keep inventing the wheel, over and over again, trying to create API for our service/product and spending enormous time and money trying to integrate different products, sometimes inside the same group... Most SE consider the infra components, like Kafka, NATS, DB, ETCD & etc., as the "Wheel" and rushing off to implement and usage those infra components, thinking they are using the "Wheel"... But in fact, **they are just doing the complete opposite.** While creating those from scratch is a nice challenge, it isn't as expensive as maintaining API and integrations over time. Putting to usage infra components is a very easy task that can take a month or even weeks, while building a stable API and integrating with different products might take years and constant costly maintenance over time.
If we do an analogy to Language, the infra components are the alphabet, while the API is the actual Languages. The same as two persons, each knows a different language but with the same alphabet, cannot speak to each other, two products, built with the same infra cannot communicate with each other and require a very expensive, highly maintenance integration.

The Graph SQL comes to ease the language/api challenge by presenting a single, simple & common API to query the model & data or a product at runtime.

## So how does it work?
### Model Inspector**
First you got the Model Inspector, the model inspector is accepting a GO struct or a Java Object and starts to introspect the struct/object, drilling down to discover its attributes and sub objects. From this data, it is creating the GSchema.

### Parser
The parser just parses the query and validates that the syntax is correct. It divides the query string to requested GColumn, GTables & Criteria. The Criteria is divided into Expression & Compares.

### InstanceId
An InstanceId is a string representation of an instance inside the model. For example if we have the following model:
   Employee
        Addresses
            [] Address
                  Line1
                  Line2
                  Zip
                    Country
To refer to an Employee instance, will do it via “Employee[key]”.
To refer to an Address instance, will do it vis “Employee[key].Adresses.Address[0..]  
