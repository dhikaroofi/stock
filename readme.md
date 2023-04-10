# Stockbit

**How To Run:**
- clone the repository “https://github.com/dhikaroofi/stock” or extract zip file
- enter to the project folder
- adjust the config file
    - copy “config.yaml.example ” and rename it to “config.yaml”
    - make sure that the configuration exactly matches your environment,
    - specially for the redis config
- using Makefile
    - **make run part1** → is to run server for challenge part 1
    - **make run part2** → is to run code for challenge part 2


**How To Code**

**The Arc**

This application is based on the hexagonal architecture, which involves separating the input/output layer from the business logic layer, also known as the "core" layer. In this architecture, the core layer is responsible for the application's business logic and domain model, while the input/output layer is responsible for communicating with the external world, such as through APIs or UIs. This separation of concerns helps to keep the core layer decoupled from external dependencies, making it easier to test and maintain over time.

**Project Structure**

- Internal
    - adapters
        - driven
          these adapters are used to interact with external systems such as databases, message queues, or file systems.
        - driving
          these adapters are used to expose the application's core logic to the outside world, such as through APIs or user interfaces.
    - cmd
      this folder is to aggregate the adapters and the core
    - config
      this component contains configuration files for the application, such as database connection settings, or environment variables.
    - core
      this is the heart of the application, containing the business logic and domain model. It is designed to be independent of external dependencies, making it easier to test and maintain.
    - entity/domain
      this folder is contain the struct or data model for the app
- resources
    - this filled with the all of the resources
- pkg

  this contain third party package

- tools

  this for help to generate the mock

- main.go

  this the entry point of the app