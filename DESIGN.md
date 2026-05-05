# Add support for lock statements

+ We have a gurantee that there is only one mutable variable referred within a lock statement
  > Only one such variable can occur in the lock statement.
  + We should use a lock id based on this variable 
    + We can introduce some new analysis env for concurrency that maps a symbol ref to lock id 
      + TODO: is this going to work even when we read the symbol from the disk consistently
    + Then we pass the id to the runtime in some sort of lock start and lock end instructions
      + Runtime needs to maintain a map for id to concurrent lock and lock that?
  + We need a new analyzer for the lock which will use the visitor to get all the references to mutable variables (similar to isolation analysis) and validate they all refer to the same variable
    + After that (may be) add the symbol ref to the mutable variable to the AST
    + Refactor isolation analysis to ignore lock statements

### Lock analyzer
1. There can be only one "mutable" variable referred, we call this the "restricted" variable.
  + Pick the first mutable variable you come across as the restricted variable.
2. Expression of `return` statement must be isolated.
3. All assignments to variables that are not defined in the lock statement and not the restricted variable must be isolated
4. All variable references to variables not defined in lock statement and not the restricted variable must happen within isolated expressions.


- I think the most straight forward way to implement this analysis is to have a lock analyzer that determine the restricted variable and give an error if there is more than one.
  - We can extract the logic in isolation analysis used to determine if the variable is isolated or not for this and reuse it in this analyzer.
- Then we should refactor isolation analysis such that we can give it additional variables as constant and pass in the restricted variable as constant within the context of statement block and run the isolation analysis within the block

### Restricted variable
- Restricted in addition to being a mutable variable must be one of the followings
  1. Module level isolated variable
    ```ballerina
      isolated int foo = 5;

      public function main() {
          bar();
      }

      isolated function bar() {
          lock {
              foo = 10;
          }
      }
    ```
    - Note: accessing module level isolated variable outside of a lock statement is an error
      ```ballerina
        isolated int foo = 5;

        public function main() {
            foo = 10; // @error
        }
        ```
  2. non final field within an isolated class
    ```ballerina
      isolated class MyIsolatedClass {
          private int foo = 10;
          final int bar = 15;

          function f() {
              lock {
                  int a = self.foo;
              }
          }
      }
    ```
    - Note non final fields within an isolated class must be private
      ```ballerina

        isolated class MyIsolatedClass {
          int foo = 10; // @error
        }
      ```


## Runtime implementation

### BIR gen
+ We will have 2 new instructions lock start and lock end both giving a lock id (`int`) 
  + both of these instructions are terminating instructions (this is prevent lock body from being in BBs outside of lock body)

### Runtime 
+ IMPORTANT: we can't just use go mutex to implement locks because Ballerina locks must be re-entrant
  > A naive implementation can simply acquire a single, program-wide, recursive mutex before executing a lock statement, and release the mutex after completing the execution of the lock statement. A more sophisticated implementation can perform compile-time analysis to infer a more fine-grained locking strategy that will have the same effects as the naive implementation.

+ So each strand is going to need an identity (currently all programs run in a single strand)
  + Each strand has it's own context so we can keep something like `StrandId Int64`. When creating the context we can use a `atomic.Int64`
    + IMPORTANT: we should wrap on overflow
+ We need to maintain global runtime stage (currently only global state is registry). For this we'll introduce Environment and move registry to Environment
  + Context will have reference to Environment
+ In the environment we will have a thread safe (behind locks) map from lock id to `ReentrantMutext`
  ```go

  // This is just a draft; feel free to improve/fix this
  type ReentrantMutex struct {
      mu    sync.Mutex
      cond  *sync.Cond
      owner StrandId
      count int
  }

  func (r *ReentrantMutex) Lock(strandId StrandId) {
      r.mu.Lock()
      defer r.mu.Unlock()
      for r.owner != nil && r.owner != strandId {
          r.cond.Wait()
      }
      r.owner = token
      r.count++
  }

  func (r *ReentrantMutex) Unlock(strandId StrandId) {
      r.mu.Lock()
      defer r.mu.Unlock()
      if r.owner != strandId {
          panic("unlock by non-owner")
      }
      r.count--
      if r.count == 0 {
          r.owner = nil
          r.cond.Signal()  // wake one waiter
      }
  }
  ```

+ Validate this by adding a recursive function call that use lock


