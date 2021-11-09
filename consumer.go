package consumer

import (
	"fmt"
	"net/http"

	"github.com/Capucinoxx/consumer/middleware"
)

// Method est une méthode du protocole
// http ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"]
type Method string

// Route est la représentation d'une route http
type Route struct {
	// description de la route
	Name string

	// méthode du protocole http ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"]
	Method Method

	// pattern suivant le préfix
	Pattern string

	HandlerFunc http.HandlerFunc
}

// Routes représente une collection de routes
type Routes []Route

// router représentation internet du router.
// À cette struct sera greffé les différentes fonctionnalités du corsp de la logique
// métier permettant la gestion de la consommation et mise en place des routes du service
type router struct {
	// ensemble des routes du service
	// map[pattern]map[method]Route
	routes map[string]map[Method]Route

	// ensemble des middlewares du service
	middlewares middleware.Middlewares

	// fonction permettant l'affichage des différentes routes lors de la consommation des
	// différentes routes
	printRoute func(pattern string, methods ...string)

	// état de l'instance si les routes ont été consummées ou non
	isConsummed bool
}

func PrintMethod(pattern string, methods ...string) {
	for i := 0; i < len(methods); i++ {
		fmt.Printf("[%-6s] %s\n", methods[i], pattern)
	}
}

// Router construit une instance de router en assigant les routes en paramèetres à la
// structure interne.
//
// De base, le Router ne comporte aucune routeet un seul middleware, soit
// middleware.Logger permettant de calculer et d'afficher le temps pris pour
// exécuter la requête utilisateur
func Router(routes Routes, middlewares ...middleware.Middleware) *router {
	// ajout du middleware Logger au début
	middlewares = append(middlewares, nil)
	copy(middlewares[1:], middlewares)
	middlewares[0] = middleware.Logger

	r := &router{
		make(map[string]map[Method]Route),
		middlewares,
		PrintMethod,
		false,
	}

	// création des différentes routes pouvant être présente
	r.makeRoutes(routes...)

	return r
}

// makeRoutes fonction interne ecapsulant la logique métier concernant l'ajout de route dans
// la représentation
func (r *router) makeRoutes(routes ...Route) {
	for i := 0; i < len(routes); i++ {
		if _, ok := r.routes[routes[i].Pattern]; !ok {
			r.routes[routes[i].Pattern] = make(map[Method]Route)
		}
		r.routes[routes[i].Pattern][routes[i].Method] = routes[i]
	}
}

// AddRoutes ajoute des routes à la liste de routes allant être utilisé
// par le service
func (r *router) AddRouter(routes ...Route) *router {
	r.makeRoutes(routes...)
	return r
}

// DeleteByMethod retire pour toutes routes présentes la méthode passé en paramètre
func (r *router) DeleteByMethod(method Method) *router {
	for _, methods := range r.routes {
		delete(methods, method)
	}

	return r
}

// DeleteByPattern retires toutes méthodes présentes pour le pattern ciblé
func (r *router) DeleteByPattern(pattern string) *router {
	delete(r.routes, pattern)
	return r
}

// AddMiddlewares ajout d'un ou plusieurs middlewares au router
func (r *router) AddMiddlewares(middlewares ...middleware.Middleware) *router {
	r.middlewares = append(r.middlewares, middlewares...)

	return r
}

// SetPrintMethod permet de customiser la méthode d'impression des routes lors de la
// consommation des routes
func (r *router) SetPrintMethod(printMethod func(pattern string, methods ...string)) *router {
	r.printRoute = printMethod

	return r
}

// WithoutLogger retire de la liste des middlewares utilisés le Logger. Comme ce dernier
// est utilisé par défault, utiliser cette fonction pour retirer le Logger.
func (r *router) WithoutLogger() *router {
	r.middlewares = r.middlewares[1:]

	return r
}

// Consumer consume les routes présentes dans la structure interne pour les implanter avec le package http
// cette méthode est finale sur la structure, les modifications futures faite sur l'instance de la structure
// router ne pourra être consumé une seconde fois.
func (r *router) Consumer(prefix string) {
	if r.isConsummed {
		panic("error")
	}
	r.isConsummed = true

	for pattern, methods := range r.routes {
		func(pattern string, methods map[Method]Route) {
			handler := func(w http.ResponseWriter, rq *http.Request) {
				if next, ok := methods[Method(rq.Method)]; ok {
					// s'il n'y a pas de middleware, on retourne la fonction
					if len(r.middlewares) == 0 {
						next.HandlerFunc(w, rq)
						return
					}

					// sinon on fait une construction inverse des middlewares pour faire une imbrication
					// ex: m1(m2(m3(next.HandlerFunc(w, rq))))
					// pour m1, m2 et m3 des middlewares
					wrapped := next.HandlerFunc
					for i := len(r.middlewares) - 1; i >= 0; i-- {
						wrapped = r.middlewares[i](wrapped)
					}
					wrapped(w, rq)
				}
			}

			http.HandleFunc(prefix+pattern, handler)

			m := make([]string, len(methods)-1)
			for method := range methods {
				m = append(m, string(method))
			}
			r.printRoute(prefix+pattern, m...)
		}(prefix+pattern, methods)
	}
}
