# Gliik: A helper spirit for the everyday

Gliik leverages AI to solve personal or collective problems across multiple contexts.

Gliik FOLLOWS THE UNIX PHILOSOPHY and has a strong emphasis on composability, minimalism, and clear separation of concerns.

- Universal access to AI capabilities through familiar interfaces
- Seamless workflow integration
- Modular design allowing piping or chaining between tools and modules
- Consistent experience across platforms

Gliik should be accesible in the most natural context of use for each case. It could be used from a browser, as an option in a email app, from a launcher extension, or a code editor, a terminal, etc. To achieve this, the core logic of the system has to take different forms: a binary CLI, a web service, even a client-side code if needed.

## Core concepts

To be capable of solving problems easily, Gliik uses different concepts in specific relations.

For example, a user tells Gliik to make an optimized version of her resume that fits better with a specific job description.To do so, she needs to use the `enhance resume` action, and provide additional context: the job offer and the base resume (both as text strings, files or URL). The action executes an instruction (a prompt sent to an AI model with the additional context) and returns the new cv.

So the concepts needed to understand the example are the following.

### Instructions

Instructions are stable prompts that tells an AI Model what to do, how to do it and what to return. Instructions solve very specific tasks in specific ways.

For example, an `enhance_resume` instruction tells the model that it will receive a job offer, a base resume, and with that it needs to analyze both, make an enhancement plan, execute it and return a new resume content.

Instructions should concern only of input/output text. Therfore, they might be capable of piping one another in the UNIX style.

### Situational context

Situational context is all the extra information required to execute an action/instruction. It can be also extra information that could improve the result, but isn't required.

Situational context should be provided explicitly by the user. But it could be inferred by extra data that Gliik can access. For example, the user preferences or user profile files, the session in which Gliik was called (which may include data about the client, device, or, if it has permissions, the date and time, location, etcetera).

To be capable of access to personal context information, Gliik needs explicit permission of the user, stored in the user config file.

### Actions

Actions execute one or more instructions, providing situational context and returning the final output to the user. Actions are the entry point for the user to access Gliik. The user can list, add, edit or remove actions.

An action can be as simple as in the `enhace resume` example, which execute the `enhance_resume` instruction with context and return the text to a stdout (note how the action has spaces in its name, but the instruction has `_`).

Also an action can be as complex as `summarize book` with extra options `personalized` and `write doc`, which may include operations such as:

- Convert pdf to plaint text
- Manage the context augmentation mechanism
- Consider the current interests of the user via an instruction that access to a personal file
- Call the summarization instruction giving prominence to the personal preferences context
- Handle errors across the pipeline of instructions
- Return the summary and convert it to a .docx file

Therefore, actions are like programmable agents. They don't have full autonomy, and benefit from user feedback depending on the complexity of the task.

## Core components

### Prompt management

Technically, instructions are folders in a separate repository named after the task they perform. They also have inside a `system.md` with detailed system instructions for the AI model. Example, an `assess_paper/system.md` that tells the model how to evaluate research papers. Instructions uses git versioning to control changes.

Repository folders can be local directories or remote ones installed via git.

Instructions repositories need an internal mechanism to track prompt versioning, it should be human readable and machine accesible. It could be a separate yaml file with fields like: `version (str), description (str)`. The prompt management Gliik component must be capable of access and update this kind of system.

### User preferences

An environment-agnostic way to know user preferences, for example, which AI models and providers available should be used, or explicit permissions for gliik. It should be a single yaml file that Gliik should read to locate dependencies like APIs, Model Context Protocols and other dependencies. Also, to point to personal additional context files (like a system file for Gliik or a user file for personalization).

All secrets, API keys and other private info must be handled by environment variables (local) or a secret service (web).

### Distribution

The capacity for the core module to be served in multiple ways: as a binary cli, as a web service, or even as a module.

This might be the most complex requirement of Gliik, since it needs a strategy for handling dependencies across different environments; or define a way to allow for customization of the pipeline based on deployment context.

Some rough ideas about this:

- Monorepo with Go Modules: Use a monorepo structure with Go modules to manage shared core logic and separate main packages for the CLI, web service, and potentially a Go-Wasm module for client-side (if feasible for the intended "client-side code").
- Docker for Web Service: Package the web service in a Docker container to encapsulate its dependencies, ensuring consistent deployment across environments.
- Conditional Compilation/Build Tags: Use Go's build tags to conditionally include or exclude code based on the target environment (e.g., cli tag for CLI-specific features, web for web service).
- Plugin/Extension System: For certain functionalities (like file conversions or integration with specific AI providers), consider a plugin-based architecture. This allows different implementations to load only the necessary plugins, reducing the overall dependency footprint for each distribution.

