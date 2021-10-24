package link

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/VolticFroogo/Froogo/db"
	"github.com/VolticFroogo/Froogo/helpers"
	"github.com/gorilla/context"
	openid "github.com/yohcop/openid-go"
)

var (
	nonceStore     = openid.NewSim