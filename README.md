# gsql
GSQL - Graph SQL, An alteration and facade of the SQL language for querying Graph Models.


## Overview
We keep inventing the wheel, over and over again, trying to create API for our service/product and spending enormous time and money trying to integrate different products, sometimes inside the same group... Most SE consider the infra components, like Kafka, NATS, DB, ETCD & etc., as the "Wheel" and rushing off to implement and usage those infra components, thinking they are using the "Wheel"... But in fact, **they are just doing the complete opposite.** While creating those from scratch is a nice challenge, it isn't as expensive as maintaining API and integrations over time. Putting to usage infra components is a very easy task that can take a month or even weeks, while building a stable **API and integrating with different products might take years, huge amount of money and constant costly maintenance over time.**

If we do an analogy to Language, the infra components are the alphabet, while the API is the actual Languages. The same as two persons, each knows a different language but with the same alphabet, cannot speak to each other, two products, built with the same infra cannot communicate with each other and require a very expensive, highly maintenance integration.

The Graph SQL comes to ease the language/api challenge by presenting a single, simple & common API to query the graph model & data of a product at runtime.

The GSQL Approach


## So how does it work?
### Introspector
First you got the Model Introspector, the model inspector is accepting a GO struct or a Java Object and starts to introspect the struct/object, drilling down to discover its attributes and sub objects. From this data, it is creating the Internal Schema, mapping a struct->table and attribute->column and creating the Graph Schema of the struct, mapping the relations between the root struct and its sub structs. The model introspector allows you to define annotations for an attribute so later on, when/if you would like to use one of the pluggable persistency or even implement your own, you can do this easiely.

**Note: The actual db functionality & persistency layers have been extracted to another project so you could use the gsql over your model without model just for sorting and filtering you model element lists.**

### Parser
The parser just parses the query and validates that the syntax is correct. It divides the query string to requested Column, Tables & Criteria. The Criteria is divided into Expression, Compares & Conditions.

### Instance
An Instance is a string representation of an instance inside the model. For example if we have the following model:

    Employee
        Addresses
            [] Address
                   Line1
                   Line2
                   Zip
                   Country


To refer to an Employee instance, will do it via **“Employee[key]”**.

To refer to an Address instance, will do it via **“Employee[key].Adresses.Address[0..]"**.  

### Attribute
An Attribute is a string representations of a struct attribute in the model, for example:

To refer to Line2 in an Instance of Address, will do it via **"Employee[key].Adresses.Address[0..].Line2"**

The Attribute also contains a seemless setter & getter from your model so you don't need worry about the instance chain/tree extracting data or setting data in your model. For example the following code will set the Line2 value even if Addresses is nil & Address list is nil:

    employee := &Employee{}
    attribute,_ := CreateAttribute("Employee.Adresses.Address[1].Line2")
    attribute.SetValue("My Address Line 2")
    
### Interpreter
The Inerpreter is taking a syntax valid parsed query and validating it via the Introspector schema, trying to match the string representation of the attributes to the discovered Columns & tables by the Introspector. If successful, the outcome is an Interperter Query instance that you can use to filter elements in your model using the Match method, e.g. the use it like the following code:

    query,_:=NewQuery("select * from mymodel where name='xxx' or family='yyy'")
    for _,modelElement:=range myModelElementsList {
        if query.Match(modelElement) {
            //This model element match the query criteria.
        }
    }
    
