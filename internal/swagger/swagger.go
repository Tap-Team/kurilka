package swagger

import (
	"github.com/Tap-Team/kurilka/docs"
	"github.com/Tap-Team/kurilka/internal/config"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

//	@title			Kurilka API Swagger
//	@version		1.0
//	@termsOfService	http://swagger.io/terms/
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
//	@BasePath		/
func Swagger(
	r *mux.Router,
	cnf config.ServerConfig,
) {
	docs.SwaggerInfo.Host = cnf.SwaggerHost
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}
