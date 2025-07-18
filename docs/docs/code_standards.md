# Code standards

## Logging rules

We came up with these rules to help debugging but also avoid clutter and info leaking to the public.

1. Always use `WrapErrorWithDetails` at the deepest level of your code. If you got an error from a code you called
   yourself, reconsider wrapping the error with this.
    - If you want to wrap an error, use our custom `ErrorWrapper`.
2. It's okay to log on every level of the stack using `LogDetailedError`.
3. Do not show the client the errors as it is! Use `ErrorWrapper.UserMessage`!

## Go imports

Follow this structure:

```go
import (
	// system
	"context"
	"errors"
	"net/http"

	// external
	"github.com/twitchtv/twirp"
	"golang.org/x/crypto/bcrypt"
	
	// internal remote: use go get
	"github.com/CherryXenon/cherry_api/rpc/cherry"

	// internal local
	"github.com/CherryXenon/Cherry/internal/server"
	"github.com/CherryXenon/Cherry/internal/util"
	"github.com/CherryXenon/Cherry/pkg/models"
)
```