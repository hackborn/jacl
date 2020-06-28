# jacl
A comparision library for Go.

### DESCRIPTION ###

Jacl is a JSON asymmetric comparison library: It uses JSON to reduce objects to structures that can be compared, and does an asymmetric comparison where all the values of A must be present in B and match, but all the values of B do not need to be present or match A.

The intended purpose is to use for testing libraries, where you often need to only check for specific values in complex result sets. Since testing is the intended use case performance has not been a consideraion.

### PRERELEASE ###

Ignore this, currently in a prerelease state.