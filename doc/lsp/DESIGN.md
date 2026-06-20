# Design

- State of each project (either single file or package) is represented by a `Snapshot` consisting of
  - `CompilerEnv`
  - topologically sorted pkg list
    - This is just index to an pkg array
  - Per module/package
    - AST
    - state
    - exported symbol space
    - per file in module
      - source
      - CU
  - Generation (int)
- When adding a project to server read all the project files from the disk and build the Snapshot with all the information (IE you need to run upto diagnostics level)
- Server will maintain multiple snapshots per each open project and direct each request with the correct snapshot
  - Use a map key by either project root or file path for single file
- We will split each message into 2 kinds
  1. Updates/notifications: 
    - These messages update Snapshot from one state to another
      - We do the minimal amount of update while invalidating others
    - These are always handled sequentially
      - **This needs to be so because depends on the other**. For example consider you get 2 edits for same file second edit assumes first one is already done.
      - *I think we can technically do these in parallel when updates are for different files but that is sufficiently rare and implementation complex enough probably not worth it at the moment*
    - **These always bump the generation of the snapshot**
  2. requests:
    - These use and update the current snapshot without bumping generation
    - These are always handled sequentially
      - No snapshot copy is needed for request handling
      - Request mutations are only cached frontend state, not source changes
- Each message we receive we add them to a queue as we finish parsing and message handler will pull messages from that queue

## Updates

### DidOpen
1. Discover the project for the opened file.
2. If the project is not tracked, add it to the server.
3. Update the file source from the request text.
   - If request text is not available, read the file from disk.
4. Make that file's CU nil.
5. Reset package state for the package containing that file.
6. Reset state for all dependent packages after this package in topsort.
7. Reset topsort.

### Project discovery
- Given a file path, walk parent directories until a `Ballerina.toml` is found.
  - If found, the project key is that directory path.
  - Add/read all project files under that root.
- If no `Ballerina.toml` is found, treat the file as a single-file project.
  - The project key is the file path.
- `initialize` discovers and adds the root project if the root is a Ballerina project.
- `didOpen` discovers and adds the project for the opened file if it is not already tracked.

### DidChange
1. Update that files source string
2. Make that CU nil
3. Reset package state for the package containing that file
4. Reset state, except source strings, for all dependent packages after this package in topsort
5. Reset topsort

### DidSave
1. Reset the file similar to Edit and reload the file content from disk.
2. Ignore request text if present.

### DidClose
- No-op.

## Request
- For handling a request we need to run certain actions against the front end to get it to a certain level after which we can do some analysis
- Dispatching an action assumes its prerequisite actions have already been dispatched. The caller dispatching the action is responsible for dispatching prerequisites.
- Each action is an atomic update on the current snapshot
  - If the snapshot has already completed the requested action, return it without redoing the action
  - We'll have a dispatch function that takes snapshot + action + *CompilerContext (some actions don't need this so nil) and return snapshot
  - Helper functions that build array of actions
- Actions:
  - `parse(file)`
    - Take the string for the file and create the CU for that file
  - `topoSort`
    - Update topo sorted pkg list in snapshot
  - `symbolResolve(pkg)` -> upto this pkg
    - IMPORTANT: you will need to inject the langlib exported symbol spaces
    - Set the snapshot module AST and update module state
  - `topLevelTypeResolve(pkg)` -> upto this pkg
    - Use snapshot module AST and update module state
  - `localTypeResolve(pkg)`
    - Use snapshot module AST and update module state
  - `semanticAnalysis(pkg)`
    - Use snapshot module AST and update module state
  - `buildCFG(pkg)`
    - Use snapshot module AST and update module state
  - `cfgAnalysis(pkg)`
    - Use snapshot module AST and update module state

### publish Diagnostics
- Actions
  1. first set it dispatched as
    - parse all files
    - toposort
    - symbolResolve (last pkg in topsort)
    - topLevelTypeResolve(last pkg in topsort)
  2. For each package, concurrently dispatch this sequential action chain
    - localTypeResolve
    - semanticAnalysis
    - buildCFG
    - cfgAnalysis
- At each action check if the compiler context has diagnostics; if so stop continuing the rest of the action set
- Accumulate diagnostics and publish them for all files in the project
  - Publish empty diagnostics for files without diagnostics to clear stale client diagnostics
  - If the snapshot generation changed before publishing, publish nothing
  - This needs to map our positions to LSP positions and add other metadata LSP needs

### Completion at point
1. Dispatch parse for all project files.
2. Dispatch topsort.
3. Check if completion is at `[text]:$`.
  - Else return nothing for now.
4. Using CU imports identify what is the module.
  - Use `[text]` to narrow down the import node. *This should handle alias*.
5. Dispatch symbolResolve up to that module.
  - If compiler context has diagnostics return nothing.
6. Dispatch topLevelTypeResolve up to that module.
  - If compiler context has diagnostics return nothing.
7. Get the exported symbol space for that package and build the completion list using both symbol name and symbol type.
