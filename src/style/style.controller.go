package style

// type patternController struct {
// 	service service.Default[m.Pattern, m.PatternQueryDto]
// }

// func NewPatternController(s service.Default[m.Pattern, m.PatternQueryDto]) controller.Default {
// 	return &patternController{service: s}
// }

// // @Tags Pattern
// // @Summary Create Pattern by project id
// // @Accept json
// // @Produce application/json
// // @Produce application/xml
// // @Security BearerAuth
// // @Param model body m.PatternDto true "Pattern info"
// // @Success 201 {object} m.Success{result=[]m.Pattern}
// // @failure 400 {object} m.Error
// // @failure 401 {object} m.Error
// // @failure 422 {object} m.Error
// // @failure 429 {object} m.Error
// // @failure 500 {object} m.Error
// // @Router /pattern [post]
// func (o *patternController) CreateOne(c *gin.Context) {
// 	var body m.PatternDto
// 	if err := c.ShouldBind(&body); err != nil || !body.IsOK() {
// 		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { body: %t }", body.IsOK()))
// 		return
// 	}

// 	var model = m.Pattern{Mode: body.Mode, Colors: body.Colors, MaxStroke: body.MaxStroke, MaxScale: body.MaxScale, MaxSpacingX: body.MaxSpacingX, MaxSpacingY: body.MaxSpacingY, Width: body.Width, Height: body.Height, Path: body.Path}
// 	if err := o.service.Create(&model); err != nil {
// 		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	helper.ResHandler(c, http.StatusCreated, &m.Success{
// 		Status: "OK",
// 		Result: []m.Pattern{model},
// 		Items:  1,
// 	})
// }

// // @Tags Pattern
// // @Summary Create Pattern list of objects
// // @Accept json
// // @Produce application/json
// // @Produce application/xml
// // @Security BearerAuth
// // @Param model body []m.PatternDto true "List of Patterns info"
// // @Success 201 {object} m.Success{result=[]m.Pattern}
// // @failure 400 {object} m.Error
// // @failure 401 {object} m.Error
// // @failure 422 {object} m.Error
// // @failure 429 {object} m.Error
// // @failure 500 {object} m.Error
// // @Router /pattern/list [post]
// func (o *patternController) CreateAll(c *gin.Context) {
// 	var body []m.PatternDto
// 	if err := c.ShouldBind(&body); err != nil || len(body) == 0 {
// 		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { body: %t }", len(body) != 0))
// 		return
// 	}

// 	var models = []m.Pattern{}
// 	for _, item := range body {
// 		if item.IsOK() {

// 			var model = m.Pattern{Mode: item.Mode, Colors: item.Colors, MaxStroke: item.MaxStroke, MaxScale: item.MaxScale, MaxSpacingX: item.MaxSpacingX, MaxSpacingY: item.MaxSpacingY, Width: item.Width, Height: item.Height, Path: item.Path}
// 			if err := o.service.Create(&model); err != nil {
// 				helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
// 				return
// 			}

// 			models = append(models, model)
// 		}
// 	}

// 	helper.ResHandler(c, http.StatusCreated, &m.Success{
// 		Status: "OK",
// 		Result: models,
// 		Items:  len(models),
// 	})
// }

// // @Tags Pattern
// // @Summary Read Pattern by :id
// // @Accept json
// // @Produce application/json
// // @Produce application/xml
// // @Param id path int true "Instance id"
// // @Success 200 {object} m.Success{result=[]m.Pattern}
// // @failure 429 {object} m.Error
// // @failure 400 {object} m.Error
// // @failure 500 {object} m.Error
// // @Router /pattern/{id} [get]
// func (o *patternController) ReadOne(c *gin.Context) {
// 	var id = helper.GetID(c, "id")

// 	if id == 0 {
// 		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { id: %t }", id != 0))
// 		return
// 	}

// 	models, err := o.service.Read(&m.PatternQueryDto{ID: uint32(id)})
// 	if err != nil {
// 		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	helper.ResHandler(c, http.StatusOK, &m.Success{
// 		Status: "OK",
// 		Result: models,
// 		Items:  len(models),
// 	})
// }

// // @Tags Pattern
// // @Summary Read Pattern by Query
// // @Accept json
// // @Produce application/json
// // @Produce application/xml
// // @Param id query int false "Type: '1'"
// // @Param mode query string false "Mode: 'fill'"
// // @Param colors query int false "colors: '2'"
// // @Param page query int false "Page: '0'"
// // @Param limit query int false "Limit: '1'"
// // @Success 200 {object} m.Success{result=[]m.Link}
// // @failure 429 {object} m.Error
// // @failure 400 {object} m.Error
// // @failure 500 {object} m.Error
// // @Router /pattern [get]
// func (o *patternController) ReadAll(c *gin.Context) {
// 	var query = m.PatternQueryDto{Page: -1}
// 	if err := c.ShouldBindQuery(&query); err != nil {
// 		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
// 		return
// 	}

// 	models, err := o.service.Read(&query)
// 	if err != nil {
// 		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	helper.ResHandler(c, http.StatusOK, &m.Success{
// 		Status: "OK",
// 		Result: models,
// 		Page:   query.Page,
// 		Limit:  query.Limit,
// 		Items:  len(models),
// 	})
// }

// // @Tags Pattern
// // @Summary Update Pattern by :id
// // @Accept json
// // @Produce application/json
// // @Produce application/xml
// // @Security BearerAuth
// // @Param id path int true "Instance id"
// // @Param model body m.PatternDto true "Pattern Data"
// // @Success 200 {object} m.Success{result=[]m.Pattern}
// // @failure 400 {object} m.Error
// // @failure 401 {object} m.Error
// // @failure 422 {object} m.Error
// // @failure 429 {object} m.Error
// // @failure 500 {object} m.Error
// // @Router /pattern/{id} [put]
// func (o *patternController) UpdateOne(c *gin.Context) {
// 	var body m.PatternDto
// 	var id = helper.GetID(c, "id")

// 	if err := c.ShouldBind(&body); err != nil || id == 0 {
// 		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { id: %t, err: %v }", id != 0, err))
// 		return
// 	}

// 	models, err := o.service.Update(&m.PatternQueryDto{ID: uint32(id)}, &m.Pattern{Mode: body.Mode, Colors: body.Colors, MaxStroke: body.MaxStroke, MaxScale: body.MaxScale, MaxSpacingX: body.MaxSpacingX, MaxSpacingY: body.MaxSpacingY, Width: body.Width, Height: body.Height, Path: body.Path})
// 	if err != nil {
// 		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	helper.ResHandler(c, http.StatusCreated, &m.Success{
// 		Status: "OK",
// 		Result: models,
// 		Items:  len(models),
// 	})
// }

// // @Tags Pattern
// // @Summary Update Pattern by Query
// // @Accept json
// // @Produce application/json
// // @Produce application/xml
// // @Security BearerAuth
// // @Param id query int false "Type: '1'"
// // @Param mode query string false "Mode: 'fill'"
// // @Param colors query int false "colors: '2'"
// // @Param project_id query string false "ProjectID: '1'"
// // @Param model body m.PatternDto true "Pattern Data"
// // @Success 200 {object} m.Success{result=[]m.Pattern}
// // @failure 400 {object} m.Error
// // @failure 401 {object} m.Error
// // @failure 422 {object} m.Error
// // @failure 429 {object} m.Error
// // @failure 500 {object} m.Error
// // @Router /pattern [put]
// func (o *patternController) UpdateAll(c *gin.Context) {
// 	var query = m.PatternQueryDto{Page: -1}
// 	if err := c.ShouldBindQuery(&query); err != nil {
// 		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
// 		return
// 	}

// 	var body m.PatternDto
// 	if err := c.ShouldBind(&body); err != nil {
// 		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
// 		return
// 	}

// 	models, err := o.service.Update(&query, &m.Pattern{Mode: body.Mode, Colors: body.Colors, MaxStroke: body.MaxStroke, MaxScale: body.MaxScale, MaxSpacingX: body.MaxSpacingX, MaxSpacingY: body.MaxSpacingY, Width: body.Width, Height: body.Height, Path: body.Path})
// 	if err != nil {
// 		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	helper.ResHandler(c, http.StatusCreated, &m.Success{
// 		Status: "OK",
// 		Result: models,
// 		Items:  len(models),
// 	})
// }

// // @Tags Pattern
// // @Summary Delete Link by :id
// // @Accept json
// // @Produce application/json
// // @Produce application/xml
// // @Security BearerAuth
// // @Param id path int true "Instance id"
// // @Success 200 {object} m.Success{result=[]string{}}
// // @failure 400 {object} m.Error
// // @failure 401 {object} m.Error
// // @failure 422 {object} m.Error
// // @failure 429 {object} m.Error
// // @failure 500 {object} m.Error
// // @Router /pattern/{id} [delete]
// func (o *patternController) DeleteOne(c *gin.Context) {
// 	var id = helper.GetID(c, "id")

// 	if id == 0 {
// 		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: { id: %t }", id != 0))
// 		return
// 	}

// 	items, err := o.service.Delete(&m.PatternQueryDto{ID: uint32(id)})
// 	if err != nil {
// 		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	helper.ResHandler(c, http.StatusOK, &m.Success{
// 		Status: "OK",
// 		Result: []string{},
// 		Items:  items,
// 	})
// }

// // @Tags Pattern
// // @Summary Delete Pattern by Query
// // @Accept json
// // @Produce application/json
// // @Produce application/xml
// // @Security BearerAuth
// // @Param id query int false "Instance :id"
// // @Param mode query string false "Mode: 'fill'"
// // @Param colors query int false "colors: '2'"
// // @Success 200 {object} m.Success{result=[]string{}}
// // @failure 400 {object} m.Error
// // @failure 401 {object} m.Error
// // @failure 422 {object} m.Error
// // @failure 429 {object} m.Error
// // @failure 500 {object} m.Error
// // @Router /pattern [delete]
// func (o *patternController) DeleteAll(c *gin.Context) {
// 	var query = m.PatternQueryDto{Page: -1}
// 	if err := c.ShouldBindQuery(&query); err != nil {
// 		helper.ErrHandler(c, http.StatusBadRequest, fmt.Sprintf("Bad request: %v", err))
// 		return
// 	}

// 	items, err := o.service.Delete(&query)
// 	if err != nil {
// 		helper.ErrHandler(c, http.StatusInternalServerError, err.Error())
// 		return
// 	}

// 	helper.ResHandler(c, http.StatusOK, &m.Success{
// 		Status: "OK",
// 		Result: []string{},
// 		Items:  items,
// 	})
// }
