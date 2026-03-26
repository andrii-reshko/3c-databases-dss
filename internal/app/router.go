package app

import (
	"html/template"
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
	}
	tmpl := template.Must(template.ParseFiles(files...))
	router.SetHTMLTemplate(tmpl)

	router.GET("/", c.EvaluationHandler.ShowMatrix)
	router.POST("/evaluations", c.EvaluationHandler.UpdateMatrix)

	altRoutes := router.Group("/alternatives")
	{
		altRoutes.GET("/", c.AlternativeHandler.ListAlternatives)
		altRoutes.GET("/new", c.AlternativeHandler.ShowAlternativeForm)
		altRoutes.POST("/new", c.AlternativeHandler.CreateAlternative)
		altRoutes.GET("/edit/:id", c.AlternativeHandler.ShowAlternativeForm)
		altRoutes.POST("/edit/:id", c.AlternativeHandler.UpdateAlternative)
		altRoutes.GET("/delete/:id", c.AlternativeHandler.DeleteAlternative)
	}

	criteriaRoutes := router.Group("/criteria")
	{
		criteriaRoutes.GET("/", c.CriterionHandler.ListCriteria)
		criteriaRoutes.GET("/new", c.CriterionHandler.ShowCriterionForm)
		criteriaRoutes.POST("/new", c.CriterionHandler.CreateCriterion)
		criteriaRoutes.GET("/edit/:id", c.CriterionHandler.ShowCriterionForm)
		criteriaRoutes.POST("/edit/:id", c.CriterionHandler.UpdateCriterion)
		criteriaRoutes.GET("/delete/:id", c.CriterionHandler.DeleteCriterion)
	}

	return router
}
