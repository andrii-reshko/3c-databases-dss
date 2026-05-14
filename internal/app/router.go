package app

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func SetupRouter(c *Container) *gin.Engine {
	router := gin.Default()

	getwd, _ := os.Getwd()

	files := []string{
		filepath.Join(getwd, "templates/partials/navbar.html"),
		filepath.Join(getwd, "templates/partials/head.html"),
		filepath.Join(getwd, "templates/partials/scripts.html"),
		filepath.Join(getwd, "templates/index.html"),
		filepath.Join(getwd, "templates/alternatives.html"),
		filepath.Join(getwd, "templates/criteria.html"),
		filepath.Join(getwd, "templates/alternative_form.html"),
		filepath.Join(getwd, "templates/criterion_form.html"),
		filepath.Join(getwd, "templates/error.html"),
		filepath.Join(getwd, "templates/ranking.html"),
		filepath.Join(getwd, "templates/rules.html"),
	}

	tmpl := template.New("").Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
	})
	tmpl, err := tmpl.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	router.SetHTMLTemplate(tmpl)

	router.GET("/", c.EvaluationHandler.ShowMatrix)
	router.POST("/evaluations", c.EvaluationHandler.UpdateMatrix)
	router.GET("/ranking", c.RankingHandler.ShowRanking)
	router.GET("/rules", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "rules.html", gin.H{"title": "Expert Logic"})
	})

	notImplemented := func(ctx *gin.Context) {
		ctx.HTML(http.StatusServiceUnavailable, "error.html", gin.H{
			"title":   "Not Implemented",
			"message": "This feature is not yet implemented.",
		})
	}
	router.POST("/import", notImplemented)
	router.GET("/sensitivity", notImplemented)

	altRoutes := router.Group("/alternatives")
	{
		altRoutes.GET("/", c.AlternativeHandler.ListAlternatives)
		altRoutes.GET("/new", c.AlternativeHandler.ShowAlternativeForm)
		altRoutes.POST("/new", c.AlternativeHandler.CreateAlternative)
		altRoutes.GET("/edit/:id", c.AlternativeHandler.ShowAlternativeForm)
		altRoutes.POST("/edit/:id", c.AlternativeHandler.UpdateAlternative)
		altRoutes.GET("/delete/:id", c.AlternativeHandler.DeleteAlternative)
		altRoutes.POST("/delete/:id", c.AlternativeHandler.DeleteAlternative)
	}

	criteriaRoutes := router.Group("/criteria")
	{
		criteriaRoutes.GET("/", c.CriterionHandler.ListCriteria)
		criteriaRoutes.GET("/new", c.CriterionHandler.ShowCriterionForm)
		criteriaRoutes.POST("/new", c.CriterionHandler.CreateCriterion)
		criteriaRoutes.GET("/edit/:id", c.CriterionHandler.ShowCriterionForm)
		criteriaRoutes.POST("/edit/:id", c.CriterionHandler.UpdateCriterion)
		criteriaRoutes.GET("/delete/:id", c.CriterionHandler.DeleteCriterion)
		criteriaRoutes.POST("/delete/:id", c.CriterionHandler.DeleteCriterion)
	}

	return router
}
