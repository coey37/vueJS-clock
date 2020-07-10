package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/VolticFroogo/Froogo/db"
	"github.com/VolticFroogo/Froogo/handler/recovery"
	"github.com/VolticFroogo/Froogo/handler/trade"
	"github.com/VolticFroogo/Froogo/handler/users"
	"github.com/VolticFroogo/Froogo/helpers"
	"github.com/Vol