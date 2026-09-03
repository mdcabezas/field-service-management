package api

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	authpkg "localis-backend/internal/auth"
	corehandler "localis-backend/internal/handler/core"
	customershandler "localis-backend/internal/handler/customers"
	geocodinghandler "localis-backend/internal/handler/geocoding"
	inventoryhandler "localis-backend/internal/handler/inventory"
	notificationshandler "localis-backend/internal/handler/notifications"
	mobilehandler "localis-backend/internal/handler/mobile"
	operationshandler "localis-backend/internal/handler/operations"
	partnershandler "localis-backend/internal/handler/partners"
	planninghandler "localis-backend/internal/handler/planning"
	sharedhandler "localis-backend/internal/handler/shared"
	"localis-backend/internal/middleware"
	postgrescore "localis-backend/internal/repository/postgres/core"
	postgrescustomers "localis-backend/internal/repository/postgres/customers"
	postgresgeocoding "localis-backend/internal/repository/postgres/geocoding"
	postgresinventory "localis-backend/internal/repository/postgres/inventory"
	postgresnotifications "localis-backend/internal/repository/postgres/notifications"
	postgresoperations "localis-backend/internal/repository/postgres/operations"
	postgrespartners "localis-backend/internal/repository/postgres/partners"
	postgresplanning "localis-backend/internal/repository/postgres/planning"
	postgresshared "localis-backend/internal/repository/postgres/shared"
	servicescore "localis-backend/internal/service/core"
	servicesinventory "localis-backend/internal/service/inventory"
	servicesnotifications "localis-backend/internal/service/notifications"
	servicesoperations "localis-backend/internal/service/operations"
	servicesplanning "localis-backend/internal/service/planning"
)

func NewRouter(jwtAuth *authpkg.JWTAuth, pool *pgxpool.Pool, photoStorageDir string) (*gin.Engine, func()) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logging())
	r.Use(middleware.CORS())
	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.BodyLimit(10 << 20))
	if err := r.SetTrustedProxies([]string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"}); err != nil {
		slog.Error("failed to set trusted proxies", "error", err)
	}

	revocation := authpkg.NewRevocationService()
	var cleanup func()
	userRepo := postgrescore.NewCoreUserRepo(pool)
	authHandler := authpkg.NewHandler(userRepo, jwtAuth, revocation)
	jwtMiddleware := authpkg.NewMiddleware(jwtAuth, revocation)

	// Repositories
	techRoleRepo := postgrescore.NewCoreTechRoleRepo(pool)
	auditLogRepo := postgrescore.NewPlanAuditLogRepo(pool)
	partnerRepo := postgrespartners.NewPartnerRepo(pool)
	partnerContactRepo := postgrespartners.NewPartnerContactRepo(pool)
	partnerAgreementRepo := postgrespartners.NewPartnerAgreementRepo(pool)
	partnerAgreementDocRepo := postgrespartners.NewPartnerAgreementDocRepo(pool)
	partnerAgreementFormRepo := postgrespartners.NewPartnerAgreementFormRepo(pool)
	slaRepo := postgrespartners.NewSLARepo(pool)
	customerRepo := postgrescustomers.NewCustomerRepo(pool)
	customerAddressRepo := postgrescustomers.NewCustomerAddressRepo(pool)
	customerAddressPartnerRepo := postgrescustomers.NewCustomerAddressPartnerRepo(pool)
	propertyRepo := postgrescustomers.NewPropertyRepo(pool)
	technicianRepo := postgrescustomers.NewTechnicianRepo(pool)
	techCertRepo := postgrescustomers.NewTechCertificationRepo(pool)
	materialRepo := postgresinventory.NewMaterialRepo(pool)
	toolRepo := postgresinventory.NewToolRepo(pool)
	eppItemRepo := postgresinventory.NewEPPItemRepo(pool)
	vehicleRepo := postgresinventory.NewVehicleRepo(pool)
	vehicleTypeRepo := postgresinventory.NewVehicleTypeRepo(pool)
	rentalRepo := postgresinventory.NewRentalRepo(pool)
	costRateRepo := postgresinventory.NewCostRateRepo(pool)
	maintRecordRepo := postgresinventory.NewMaintenanceRecordRepo(pool)
	maintScheduleRepo := postgresinventory.NewMaintenanceScheduleRepo(pool)
	checklistTmplRepo := postgresinventory.NewChecklistTemplateRepo(pool)
	dailyPlanRepo := postgresplanning.NewDailyPlanRepo(pool)
	dailyPlanAssignmentRepo := postgresplanning.NewDailyPlanAssignmentRepo(pool)
	dailyLoadMatRepo := postgresplanning.NewDailyLoadMaterialRepo(pool)
	dailyLoadToolRepo := postgresplanning.NewDailyLoadToolRepo(pool)
	dailyLoadEPPRepo := postgresplanning.NewDailyLoadEPPRepo(pool)
	routeRepo := postgresplanning.NewRouteRepo(pool)
	routeTypeRepo := postgresplanning.NewRouteTypeRepo(pool)
	visitRepo := postgresoperations.NewVisitRepo(pool)
	visitTypeRepo := postgresoperations.NewVisitTypeRepo(pool)
	visitAssignmentRepo := postgresoperations.NewVisitAssignmentRepo(pool)
	ctmRepo := postgresoperations.NewChecklistTemplateMaterialRepo(pool)
	cttRepo := postgresoperations.NewChecklistTemplateToolRepo(pool)
	ctaRepo := postgresoperations.NewChecklistTemplateEPPRepo(pool)
	vcmRepo := postgresoperations.NewVisitChecklistMaterialRepo(pool)
	vctRepo := postgresoperations.NewVisitChecklistToolRepo(pool)
	vceRepo := postgresoperations.NewVisitChecklistEPPRepo(pool)
	vmuRepo := postgresoperations.NewVisitMaterialUsageRepo(pool)
	vtuRepo := postgresoperations.NewVisitToolUsageRepo(pool)
	veuRepo := postgresoperations.NewVisitEPPUsageRepo(pool)
	vmRepo := postgresoperations.NewVisitMeasurementRepo(pool)
	vcpRepo := postgresoperations.NewVisitCheckpointRepo(pool)
	vphRepo := postgresoperations.NewVisitPhotoRepo(pool)
	vrRepo := postgresoperations.NewVisitReportRepo(pool)
	riRepo := postgresoperations.NewReportImageRepo(pool)
	reRepo := postgresoperations.NewReportEntryRepo(pool)
	vaRepo := postgresoperations.NewVehicleAssignmentRepo(pool)
	vrentRepo := postgresoperations.NewVisitRentalRepo(pool)
	vslaRepo := postgresoperations.NewVisitSLATrackingRepo(pool)
	notifTmplRepo := postgresnotifications.NewNotificationTemplateRepo(pool)
	notifRepo := postgresnotifications.NewNotificationRepo(pool)
	addrRepo := postgresgeocoding.NewGeocodingAddressRepo(pool)
	reportTmplRepo := postgresshared.NewReportTemplateRepo(pool)

	// Services
	auditSvc := servicescore.NewAuditService(auditLogRepo)
	notifSvc := servicesnotifications.NewNotificationService(notifRepo)
	visitSvc := servicesoperations.NewVisitService(servicesoperations.VisitServiceDeps{
		VisitRepo:      visitRepo,
		AssignmentRepo: visitAssignmentRepo,
		TechnicianRepo: technicianRepo,
		CheckpointRepo: vcpRepo,
		CTMRepo:        ctmRepo,
		CTTRepo:        cttRepo,
		CTARepo:        ctaRepo,
		VCMRepo:        vcmRepo,
		VCTRepo:        vctRepo,
		VCERepo:        vceRepo,
		AuditService:   auditSvc,
	})
	visitAssignmentSvc := servicesoperations.NewVisitAssignmentService(visitAssignmentRepo, visitRepo, technicianRepo, auditSvc)
	dailyPlanSvc := servicesplanning.NewDailyPlanService(dailyPlanRepo, dailyPlanAssignmentRepo, technicianRepo, auditSvc)
	dailyLoadSvc := servicesplanning.NewDailyLoadService(dailyLoadMatRepo, dailyLoadToolRepo, dailyLoadEPPRepo, materialRepo, toolRepo, eppItemRepo, auditSvc)
	vehicleAssignmentSvc := servicesoperations.NewVehicleAssignmentService(vaRepo, vehicleRepo, auditSvc)
	slaTrackingSvc := servicesoperations.NewSLATrackingService(vslaRepo, slaRepo, visitRepo, auditSvc, notifSvc)
	maintSvc := servicesinventory.NewMaintenanceService(maintScheduleRepo, maintRecordRepo, vehicleRepo, auditSvc, notifSvc)
	reportSvc := servicesoperations.NewReportService(vrRepo, reRepo, riRepo, vphRepo, vmRepo, reportTmplRepo, auditSvc)

	// Handlers (services injected)
	visitH := operationshandler.NewVisitHandler(visitSvc)
	visitAssignmentH := operationshandler.NewVisitAssignmentHandler(visitAssignmentSvc)
	dailyPlanH := planninghandler.NewDailyPlanHandler(dailyPlanSvc)
	dailyLoadMatH := planninghandler.NewDailyLoadMaterialHandler(dailyLoadSvc)
	dailyLoadToolH := planninghandler.NewDailyLoadToolHandler(dailyLoadSvc)
	dailyLoadEPPH := planninghandler.NewDailyLoadEPPHandler(dailyLoadSvc)
	vaH := operationshandler.NewVehicleAssignmentHandler(vehicleAssignmentSvc)
	vslaH := operationshandler.NewVisitSLATrackingHandler(slaTrackingSvc)
	maintRecordH := inventoryhandler.NewMaintenanceRecordHandler(maintSvc)
	maintScheduleH := inventoryhandler.NewMaintenanceScheduleHandler(maintSvc)
	vrH := operationshandler.NewVisitReportHandler(reportSvc)
	notifH := notificationshandler.NewNotificationHandler(notifSvc)

	// Handlers (repo-only, no business logic)
	userH := corehandler.NewUserHandler(userRepo)
	techRoleH := corehandler.NewTechRoleHandler(techRoleRepo)
	partnerH := partnershandler.NewPartnerHandler(partnerRepo)
	partnerContactH := partnershandler.NewPartnerContactHandler(partnerContactRepo)
	partnerAgreementH := partnershandler.NewPartnerAgreementHandler(partnerAgreementRepo)
	partnerAgreementDocH := partnershandler.NewPartnerAgreementDocHandler(partnerAgreementDocRepo)
	partnerAgreementFormH := partnershandler.NewPartnerAgreementFormHandler(partnerAgreementFormRepo)
	slaH := partnershandler.NewSLAHandler(slaRepo)
	customerH := customershandler.NewCustomerHandler(customerRepo)
	customerAddressH := customershandler.NewCustomerAddressHandler(customerAddressRepo)
	customerAddressPartnerH := customershandler.NewCustomerAddressPartnerHandler(customerAddressPartnerRepo)
	propertyH := customershandler.NewPropertyHandler(propertyRepo)
	technicianH := customershandler.NewTechnicianHandler(technicianRepo)
	techCertH := customershandler.NewTechCertificationHandler(techCertRepo)
	materialH := inventoryhandler.NewMaterialHandler(materialRepo)
	toolH := inventoryhandler.NewToolHandler(toolRepo)
	eppItemH := inventoryhandler.NewEPPItemHandler(eppItemRepo)
	vehicleH := inventoryhandler.NewVehicleHandler(vehicleRepo)
	vehicleTypeH := inventoryhandler.NewVehicleTypeHandler(vehicleTypeRepo)
	rentalH := inventoryhandler.NewRentalHandler(rentalRepo)
	costRateH := inventoryhandler.NewCostRateHandler(costRateRepo)
	checklistTmplH := inventoryhandler.NewChecklistTemplateHandler(checklistTmplRepo)
	routeH := planninghandler.NewRouteHandler(routeRepo)
	routeTypeH := planninghandler.NewRouteTypeHandler(routeTypeRepo)
		ctmH := operationshandler.NewChecklistTemplateMaterialHandler(ctmRepo)
		cttH := operationshandler.NewChecklistTemplateToolHandler(cttRepo)
		ctaH := operationshandler.NewChecklistTemplateEPPHandler(ctaRepo)
		checklistStageH := operationshandler.NewChecklistStageHandler()
		vcmH := operationshandler.NewVisitChecklistMaterialHandler(vcmRepo)
		vctH := operationshandler.NewVisitChecklistToolHandler(vctRepo)
		vceH := operationshandler.NewVisitChecklistEPPHandler(vceRepo)
		vmuH := operationshandler.NewVisitMaterialUsageHandler(vmuRepo)
		vtuH := operationshandler.NewVisitToolUsageHandler(vtuRepo)
		veuH := operationshandler.NewVisitEPPUsageHandler(veuRepo)
		vmH := operationshandler.NewVisitMeasurementHandler(vmRepo)
		vcpH := operationshandler.NewVisitCheckpointHandler(vcpRepo)
		vphH := operationshandler.NewVisitPhotoHandler(vphRepo)
		riH := operationshandler.NewReportImageHandler(riRepo)
		reH := operationshandler.NewReportEntryHandler(reRepo)
		vrentH := operationshandler.NewVisitRentalHandler(vrentRepo)
		notifTmplH := notificationshandler.NewNotificationTemplateHandler(notifTmplRepo)
		addrH := geocodinghandler.NewAddressHandler(addrRepo)
		reportTmplH := sharedhandler.NewReportTemplateHandler(reportTmplRepo)
		dailyPlanAssignmentH := planninghandler.NewDailyPlanAssignmentHandler(dailyPlanAssignmentRepo)
		visitTypeH := operationshandler.NewVisitTypeHandler(visitTypeRepo)

		// Mobile handlers
		photoH := mobilehandler.NewPhotoHandler(photoStorageDir, pool)
		syncH := mobilehandler.NewSyncHandler(pool)
		configH := mobilehandler.NewConfigHandler()
		gpsH := mobilehandler.NewGPSHandler()
		auditH := mobilehandler.NewAuditHandler()
		propertyHistoryH := mobilehandler.NewPropertyHistoryHandler()
		mobilePartnerH := mobilehandler.NewPartnerHandler()

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Auth
	loginRateLimiter := middleware.NewRateLimiter(500, 1*time.Minute)
	refreshRateLimiter := middleware.NewRateLimiter(30, 1*time.Minute)
	auth := r.Group("/auth")
	{
		auth.POST("/login", loginRateLimiter.Middleware(), authHandler.Login)
		auth.POST("/refresh", refreshRateLimiter.Middleware(), authHandler.Refresh)
		auth.POST("/logout", jwtMiddleware.Validate(), authHandler.Logout)
		auth.GET("/me", jwtMiddleware.Validate(), authHandler.Me)
	}

	apiRateLimiter := middleware.NewRateLimiter(1000, 1*time.Minute)
	cleanup = func() {
		revocation.Stop()
		loginRateLimiter.Stop()
		refreshRateLimiter.Stop()
		apiRateLimiter.Stop()
	}

	// Public photo serving endpoints (no auth required for direct file access)
	r.GET("/api/photos/:id/file", photoH.ServeFile)
	r.GET("/api/photos/:id/thumb", photoH.ServeThumb)

	api := r.Group("/api")
	api.Use(jwtMiddleware.Validate())
	api.Use(apiRateLimiter.Middleware())
	{
		adminOnly := authpkg.RequireRole("admin")
		adminManager := authpkg.RequireRole("admin", "manager")
		adminManagerSupervisor := authpkg.RequireRole("admin", "manager", "supervisor")

		// ===== Core =====
		api.GET("/users", adminManager, userH.List)
		api.GET("/users/:id", adminManager, userH.GetByID)
		api.POST("/users", adminOnly, userH.Create)
		api.PUT("/users/:id", adminOnly, userH.Update)
		api.DELETE("/users/:id", adminOnly, userH.Delete)

		api.GET("/tech-roles", techRoleH.List)
		api.GET("/tech-roles/:id", techRoleH.GetByID)
		api.POST("/tech-roles", adminOnly, techRoleH.Create)
		api.PUT("/tech-roles/:id", adminOnly, techRoleH.Update)
		api.DELETE("/tech-roles/:id", adminOnly, techRoleH.Delete)

		// ===== Partners =====
		api.GET("/partners", partnerH.List)
		api.GET("/partners/:id", partnerH.GetByID)
		api.POST("/partners", adminManager, partnerH.Create)
		api.PUT("/partners/:id", adminManager, partnerH.Update)
		api.DELETE("/partners/:id", adminManager, partnerH.Delete)

		partnerContacts := api.Group("/partners/:id/contacts")
		{
			partnerContacts.GET("", partnerContactH.ListByPartner)
		}
		api.GET("/partner-contacts/:id", partnerContactH.GetByID)
		api.POST("/partner-contacts", adminManager, partnerContactH.Create)
		api.DELETE("/partner-contacts/:id", adminManager, partnerContactH.Delete)

		partnerAgreements := api.Group("/partners/:id/agreements")
		{
			partnerAgreements.GET("", partnerAgreementH.ListByPartner)
		}
		api.GET("/partner-agreements/:id", partnerAgreementH.GetByID)
		api.POST("/partner-agreements", adminManager, partnerAgreementH.Create)
		api.PUT("/partner-agreements/:id", adminManager, partnerAgreementH.Update)
		api.DELETE("/partner-agreements/:id", adminManager, partnerAgreementH.Delete)

		agreementDocs := api.Group("/partner-agreements/:id/docs")
		{
			agreementDocs.GET("", partnerAgreementDocH.ListByAgreement)
		}
		api.GET("/partner-agreement-docs/:id", partnerAgreementDocH.GetByID)
		api.POST("/partner-agreement-docs", adminManager, partnerAgreementDocH.Create)
		api.DELETE("/partner-agreement-docs/:id", adminManager, partnerAgreementDocH.Delete)

		agreementForms := api.Group("/partner-agreements/:id/forms")
		{
			agreementForms.GET("", partnerAgreementFormH.ListByAgreement)
		}
		api.GET("/partner-agreement-forms/:id", partnerAgreementFormH.GetByID)
		api.POST("/partner-agreement-forms", adminManager, partnerAgreementFormH.Create)
		api.DELETE("/partner-agreement-forms/:id", adminManager, partnerAgreementFormH.Delete)

		partnerSLAs := api.Group("/partners/:id/slas")
		{
			partnerSLAs.GET("", slaH.ListByPartner)
		}
		api.GET("/slas/:id", slaH.GetByID)
		api.POST("/slas", adminManager, slaH.Create)
		api.PUT("/slas/:id", adminManager, slaH.Update)
		api.DELETE("/slas/:id", adminManager, slaH.Delete)

		// ===== Customers =====
		api.GET("/customers", customerH.List)
		api.GET("/customers/:id", customerH.GetByID)
		api.POST("/customers", adminManager, customerH.Create)
		api.PUT("/customers/:id", adminManager, customerH.Update)
		api.DELETE("/customers/:id", adminManager, customerH.Delete)

		customerAddrs := api.Group("/customers/:id/addresses")
		{
			customerAddrs.GET("", customerAddressH.ListByCustomer)
		}
		api.GET("/customer-addresses/:id", customerAddressH.GetByID)
		api.POST("/customer-addresses", adminManager, customerAddressH.Create)
		api.PUT("/customer-addresses/:id", adminManager, customerAddressH.Update)
		api.DELETE("/customer-addresses/:id", adminManager, customerAddressH.Delete)

		custAddrPartners := api.Group("/customer-addresses/:id/partners")
		{
			custAddrPartners.GET("", customerAddressPartnerH.ListByCustomerAddress)
		}
		api.GET("/customer-address-partners/:id", customerAddressPartnerH.GetByID)
		api.POST("/customer-address-partners", adminManager, customerAddressPartnerH.Create)
		api.DELETE("/customer-address-partners/:id", adminManager, customerAddressPartnerH.Delete)

		custAddrProps := api.Group("/customer-addresses/:id/properties")
		{
			custAddrProps.GET("", propertyH.ListByCustomerAddress)
		}
		api.GET("/properties/:id", propertyH.GetByID)
		api.POST("/properties", adminManager, propertyH.Create)
		api.PUT("/properties/:id", adminManager, propertyH.Update)
		api.DELETE("/properties/:id", adminManager, propertyH.Delete)

		api.GET("/technicians", technicianH.List)
		api.GET("/technicians/:id", technicianH.GetByID)
		api.POST("/technicians", adminManager, technicianH.Create)
		api.PUT("/technicians/:id", adminManager, technicianH.Update)
		api.DELETE("/technicians/:id", adminManager, technicianH.Delete)

		techCerts := api.Group("/technicians/:id/certifications")
		{
			techCerts.GET("", techCertH.ListByTech)
		}
		api.GET("/tech-certifications/:id", techCertH.GetByID)
		api.POST("/tech-certifications", adminManager, techCertH.Create)
		api.DELETE("/tech-certifications/:id", adminManager, techCertH.Delete)

		// ===== Inventory =====
		api.GET("/materials", materialH.List)
		api.GET("/materials/:id", materialH.GetByID)
		api.POST("/materials", adminManager, materialH.Create)
		api.PUT("/materials/:id", adminManager, materialH.Update)
		api.DELETE("/materials/:id", adminManager, materialH.Delete)

		api.GET("/tools", toolH.List)
		api.GET("/tools/:id", toolH.GetByID)
		api.POST("/tools", adminManager, toolH.Create)
		api.PUT("/tools/:id", adminManager, toolH.Update)
		api.DELETE("/tools/:id", adminManager, toolH.Delete)

		api.GET("/epp-items", eppItemH.List)
		api.GET("/epp-items/:id", eppItemH.GetByID)
		api.POST("/epp-items", adminManager, eppItemH.Create)
		api.PUT("/epp-items/:id", adminManager, eppItemH.Update)
		api.DELETE("/epp-items/:id", adminManager, eppItemH.Delete)

		api.GET("/vehicles", vehicleH.List)
		api.GET("/vehicles/:id", vehicleH.GetByID)
		api.POST("/vehicles", adminManager, vehicleH.Create)
		api.PUT("/vehicles/:id", adminManager, vehicleH.Update)
		api.DELETE("/vehicles/:id", adminManager, vehicleH.Delete)

		api.GET("/vehicle-types", vehicleTypeH.List)

		api.GET("/rentals", rentalH.List)
		api.GET("/rentals/:id", rentalH.GetByID)
		api.POST("/rentals", adminManager, rentalH.Create)
		api.PUT("/rentals/:id", adminManager, rentalH.Update)
		api.DELETE("/rentals/:id", adminManager, rentalH.Delete)

		api.GET("/cost-rates", costRateH.List)
		api.GET("/cost-rates/:id", costRateH.GetByID)
		api.POST("/cost-rates", adminManager, costRateH.Create)
		api.PUT("/cost-rates/:id", adminManager, costRateH.Update)
		api.DELETE("/cost-rates/:id", adminManager, costRateH.Delete)

		api.GET("/maintenance-records", maintRecordH.List)
		api.GET("/maintenance-records/:id", maintRecordH.GetByID)
		api.POST("/maintenance-records", adminManager, maintRecordH.Create)
		api.PUT("/maintenance-records/:id", adminManager, maintRecordH.Update)
		api.DELETE("/maintenance-records/:id", adminManager, maintRecordH.Delete)

		api.GET("/maintenance-schedules", maintScheduleH.List)
		api.GET("/maintenance-schedules/:id", maintScheduleH.GetByID)
		api.POST("/maintenance-schedules", adminManager, maintScheduleH.Create)
		api.PUT("/maintenance-schedules/:id", adminManager, maintScheduleH.Update)
		api.DELETE("/maintenance-schedules/:id", adminManager, maintScheduleH.Delete)

		api.GET("/checklist-templates", checklistTmplH.List)
		api.GET("/checklist-templates/:id", checklistTmplH.GetByID)
		api.POST("/checklist-templates", adminManager, checklistTmplH.Create)
		api.PUT("/checklist-templates/:id", adminManager, checklistTmplH.Update)
		api.DELETE("/checklist-templates/:id", adminManager, checklistTmplH.Delete)

		ctmRoutes := api.Group("/checklist-templates/:id/materials")
		{
			ctmRoutes.GET("", ctmH.ListByTemplate)
		}
		api.GET("/checklist-template-materials/:id", ctmH.GetByID)
		api.POST("/checklist-template-materials", adminManager, ctmH.Create)
		api.DELETE("/checklist-template-materials/:id", adminManager, ctmH.Delete)

		cttRoutes := api.Group("/checklist-templates/:id/tools")
		{
			cttRoutes.GET("", cttH.ListByTemplate)
		}
		api.GET("/checklist-template-tools/:id", cttH.GetByID)
		api.POST("/checklist-template-tools", adminManager, cttH.Create)
		api.DELETE("/checklist-template-tools/:id", adminManager, cttH.Delete)

		ctaRoutes := api.Group("/checklist-templates/:id/epps")
		{
			ctaRoutes.GET("", ctaH.ListByTemplate)
		}
		api.GET("/checklist-template-epps/:id", ctaH.GetByID)
		api.POST("/checklist-template-epps", adminManager, ctaH.Create)
		api.DELETE("/checklist-template-epps/:id", adminManager, ctaH.Delete)

		// Checklist stages
		ctStageRoutes := api.Group("/checklist-templates/:id/stages")
		{
			ctStageRoutes.GET("", checklistStageH.ListByTemplate)
		}
		api.POST("/checklist-stages", adminManager, checklistStageH.Create)
		api.PUT("/checklist-stages/:id", adminManager, checklistStageH.Update)
		api.DELETE("/checklist-stages/:id", adminManager, checklistStageH.Delete)

		// ===== Planning =====
		api.GET("/daily-plans", dailyPlanH.List)
		api.GET("/daily-plans/:id", dailyPlanH.GetByID)
		api.POST("/daily-plans", adminManager, dailyPlanH.Create)
		api.PUT("/daily-plans/:id", adminManager, dailyPlanH.Update)
		api.DELETE("/daily-plans/:id", adminManager, dailyPlanH.Delete)

		dailyPlanAssigns := api.Group("/daily-plans/:id/assignments")
		{
			dailyPlanAssigns.GET("", dailyPlanAssignmentH.ListByPlan)
		}
		api.GET("/daily-plan-assignments", dailyPlanAssignmentH.List)
		api.GET("/daily-plan-assignments/:id", dailyPlanAssignmentH.GetByID)
		api.POST("/daily-plan-assignments", adminManager, dailyPlanAssignmentH.Create)
		api.PUT("/daily-plan-assignments/:id", adminManager, dailyPlanAssignmentH.Update)
		api.DELETE("/daily-plan-assignments/:id", adminManager, dailyPlanAssignmentH.Delete)

		dailyPlanMats := api.Group("/daily-plans/:id/materials")
		{
			dailyPlanMats.GET("", dailyLoadMatH.ListByPlan)
		}
		api.GET("/daily-load-materials/:id", dailyLoadMatH.GetByID)
		api.POST("/daily-load-materials", adminManager, dailyLoadMatH.Create)
		api.PUT("/daily-load-materials/:id", adminManager, dailyLoadMatH.Update)
		api.DELETE("/daily-load-materials/:id", adminManager, dailyLoadMatH.Delete)

		dailyPlanTools := api.Group("/daily-plans/:id/tools")
		{
			dailyPlanTools.GET("", dailyLoadToolH.ListByPlan)
		}
		api.GET("/daily-load-tools/:id", dailyLoadToolH.GetByID)
		api.POST("/daily-load-tools", adminManager, dailyLoadToolH.Create)
		api.PUT("/daily-load-tools/:id", adminManager, dailyLoadToolH.Update)
		api.DELETE("/daily-load-tools/:id", adminManager, dailyLoadToolH.Delete)

		dailyPlanEPPs := api.Group("/daily-plans/:id/epps")
		{
			dailyPlanEPPs.GET("", dailyLoadEPPH.ListByPlan)
		}
		api.GET("/daily-load-epps/:id", dailyLoadEPPH.GetByID)
		api.POST("/daily-load-epps", adminManager, dailyLoadEPPH.Create)
		api.PUT("/daily-load-epps/:id", adminManager, dailyLoadEPPH.Update)
		api.DELETE("/daily-load-epps/:id", adminManager, dailyLoadEPPH.Delete)

		api.GET("/route-types", routeTypeH.List)
		api.GET("/routes", routeH.List)
		api.GET("/routes/:id", routeH.GetByID)
		api.POST("/routes", adminManager, routeH.Create)
		api.PUT("/routes/:id", adminManager, routeH.Update)
		api.DELETE("/routes/:id", adminManager, routeH.Delete)

		// ===== Operations - Visits =====
		api.GET("/visits", visitH.List)
		api.GET("/visits/:id", visitH.GetByID)
		api.POST("/visits", adminManagerSupervisor, visitH.Create)
		api.PUT("/visits/:id", adminManagerSupervisor, visitH.Update)
		api.DELETE("/visits/:id", adminManager, visitH.Delete)
		api.POST("/visits/:id/reopen", adminManagerSupervisor, visitH.Reopen)

		visitAssigns := api.Group("/visits/:id/assignments")
		{
			visitAssigns.GET("", visitAssignmentH.ListByVisit)
		}
		api.GET("/visit-assignments/:id", visitAssignmentH.GetByID)
		api.POST("/visit-assignments", adminManagerSupervisor, visitAssignmentH.Create)
		api.DELETE("/visit-assignments/:id", adminManagerSupervisor, visitAssignmentH.Delete)

		// Visit checklists
		visitCheckMats := api.Group("/visits/:id/checklist-materials")
		{
			visitCheckMats.GET("", vcmH.ListByVisit)
		}
		api.GET("/visit-checklist-materials/:id", vcmH.GetByID)
		api.POST("/visit-checklist-materials", adminManagerSupervisor, vcmH.Create)
		api.PUT("/visit-checklist-materials/:id", adminManagerSupervisor, vcmH.Update)
		api.DELETE("/visit-checklist-materials/:id", adminManagerSupervisor, vcmH.Delete)

		visitCheckTools := api.Group("/visits/:id/checklist-tools")
		{
			visitCheckTools.GET("", vctH.ListByVisit)
		}
		api.GET("/visit-checklist-tools/:id", vctH.GetByID)
		api.POST("/visit-checklist-tools", adminManagerSupervisor, vctH.Create)
		api.PUT("/visit-checklist-tools/:id", adminManagerSupervisor, vctH.Update)
		api.DELETE("/visit-checklist-tools/:id", adminManagerSupervisor, vctH.Delete)

		visitCheckEPPs := api.Group("/visits/:id/checklist-epps")
		{
			visitCheckEPPs.GET("", vceH.ListByVisit)
		}
		api.GET("/visit-checklist-epps/:id", vceH.GetByID)
		api.POST("/visit-checklist-epps", adminManagerSupervisor, vceH.Create)
		api.PUT("/visit-checklist-epps/:id", adminManagerSupervisor, vceH.Update)
		api.DELETE("/visit-checklist-epps/:id", adminManagerSupervisor, vceH.Delete)

		// Visit usages
		visitMatUsages := api.Group("/visits/:id/material-usages")
		{
			visitMatUsages.GET("", vmuH.ListByVisit)
		}
		api.GET("/visit-material-usages/:id", vmuH.GetByID)
		api.POST("/visit-material-usages", adminManagerSupervisor, vmuH.Create)
		api.DELETE("/visit-material-usages/:id", adminManagerSupervisor, vmuH.Delete)

		visitToolUsages := api.Group("/visits/:id/tool-usages")
		{
			visitToolUsages.GET("", vtuH.ListByVisit)
		}
		api.GET("/visit-tool-usages/:id", vtuH.GetByID)
		api.POST("/visit-tool-usages", adminManagerSupervisor, vtuH.Create)
		api.DELETE("/visit-tool-usages/:id", adminManagerSupervisor, vtuH.Delete)

		visitEPPUsages := api.Group("/visits/:id/epp-usages")
		{
			visitEPPUsages.GET("", veuH.ListByVisit)
		}
		api.GET("/visit-epp-usages/:id", veuH.GetByID)
		api.POST("/visit-epp-usages", adminManagerSupervisor, veuH.Create)
		api.DELETE("/visit-epp-usages/:id", adminManagerSupervisor, veuH.Delete)

		// Visit measurements
		visitMeasurements := api.Group("/visits/:id/measurements")
		{
			visitMeasurements.GET("", vmH.ListByVisit)
		}
		api.GET("/visit-measurements/:id", vmH.GetByID)
		api.POST("/visit-measurements", adminManagerSupervisor, vmH.Create)
		api.PUT("/visit-measurements/:id", adminManagerSupervisor, vmH.Update)
		api.DELETE("/visit-measurements/:id", adminManagerSupervisor, vmH.Delete)

		// Visit checkpoints
		visitCheckpoints := api.Group("/visits/:id/checkpoints")
		{
			visitCheckpoints.GET("", vcpH.ListByVisit)
		}
		api.GET("/visit-checkpoints/:id", vcpH.GetByID)
		api.POST("/visit-checkpoints", adminManagerSupervisor, vcpH.Create)
		api.DELETE("/visit-checkpoints/:id", adminManagerSupervisor, vcpH.Delete)

		// Visit photos
		visitPhotos := api.Group("/visits/:id/photos")
		{
			visitPhotos.GET("", vphH.ListByVisit)
		}
		api.GET("/visit-photos/:id", vphH.GetByID)
		api.POST("/visit-photos", adminManagerSupervisor, vphH.Create)
		api.DELETE("/visit-photos/:id", adminManagerSupervisor, vphH.Delete)

		// Visit reports
		visitReports := api.Group("/visits/:id/reports")
		{
			visitReports.GET("", vrH.ListByVisit)
		}
		api.GET("/visit-reports", vrH.List)
		api.GET("/visit-reports/:id", vrH.GetByID)
		api.POST("/visit-reports", adminManagerSupervisor, vrH.Create)
		api.DELETE("/visit-reports/:id", adminManagerSupervisor, vrH.Delete)

		// Visit types
		api.GET("/visit-types", visitTypeH.List)

		reportImages := api.Group("/visit-reports/:id/images")
		{
			reportImages.GET("", riH.ListByReport)
		}
		api.GET("/report-images/:id", riH.GetByID)
		api.POST("/report-images", adminManagerSupervisor, riH.Create)
		api.DELETE("/report-images/:id", adminManagerSupervisor, riH.Delete)

		reportEntries := api.Group("/visit-reports/:id/entries")
		{
			reportEntries.GET("", reH.ListByReport)
		}
		api.GET("/report-entries/:id", reH.GetByID)
		api.POST("/report-entries", adminManagerSupervisor, reH.Create)
		api.DELETE("/report-entries/:id", adminManagerSupervisor, reH.Delete)

		// Vehicle assignments, rentals, SLA
		api.GET("/vehicle-assignments", vaH.List)
		api.GET("/vehicle-assignments/:id", vaH.GetByID)
		api.POST("/vehicle-assignments", adminManager, vaH.Create)
		api.PUT("/vehicle-assignments/:id", adminManager, vaH.Update)
		api.DELETE("/vehicle-assignments/:id", adminManager, vaH.Delete)

		visitRentals := api.Group("/visits/:id/rentals")
		{
			visitRentals.GET("", vrentH.ListByVisit)
		}
		api.GET("/visit-rentals/:id", vrentH.GetByID)
		api.POST("/visit-rentals", adminManagerSupervisor, vrentH.Create)
		api.DELETE("/visit-rentals/:id", adminManagerSupervisor, vrentH.Delete)

		visitSLAs := api.Group("/visits/:id/sla-trackings")
		{
			visitSLAs.GET("", vslaH.ListByVisit)
		}
		api.GET("/visit-sla-trackings", vslaH.List)
		api.GET("/visit-sla-trackings/:id", vslaH.GetByID)
		api.POST("/visit-sla-trackings", adminManagerSupervisor, vslaH.Create)
		api.PUT("/visit-sla-trackings/:id", adminManagerSupervisor, vslaH.Update)
		api.DELETE("/visit-sla-trackings/:id", adminManagerSupervisor, vslaH.Delete)

		// ===== Notifications =====
		api.GET("/notification-templates", notifTmplH.List)
		api.GET("/notification-templates/:id", notifTmplH.GetByID)
		api.POST("/notification-templates", adminOnly, notifTmplH.Create)
		api.PUT("/notification-templates/:id", adminOnly, notifTmplH.Update)
		api.DELETE("/notification-templates/:id", adminOnly, notifTmplH.Delete)

		api.GET("/notifications", notifH.List)
		api.GET("/notifications/:id", notifH.GetByID)
		api.POST("/notifications", adminManager, notifH.Create)
		api.PUT("/notifications/:id", adminManager, notifH.Update)
		api.DELETE("/notifications/:id", adminManager, notifH.Delete)

		// ===== Geocoding =====
		api.GET("/addresses", addrH.List)
		api.GET("/addresses/nearby", addrH.FindNearby)
		api.GET("/addresses/:id", addrH.GetByID)
		api.POST("/addresses", adminManager, addrH.Create)
		api.PUT("/addresses/:id", adminManager, addrH.Update)
		api.DELETE("/addresses/:id", adminManager, addrH.Delete)

		// ===== Shared =====
		api.GET("/report-templates", reportTmplH.List)
		api.GET("/report-templates/:id", reportTmplH.GetByID)
		api.POST("/report-templates", adminOnly, reportTmplH.Create)
		api.PUT("/report-templates/:id", adminOnly, reportTmplH.Update)
		api.DELETE("/report-templates/:id", adminOnly, reportTmplH.Delete)

		// ===== Mobile =====
		api.POST("/mobile/photos", photoH.Upload)
		api.POST("/mobile/sync", syncH.Push)
		api.GET("/mobile/sync", syncH.Pull)
		api.POST("/mobile/sync/retry", syncH.Retry)
		api.GET("/mobile/config", configH.GetConfig)
		api.PUT("/mobile/config", configH.UpdateConfig)
		api.POST("/mobile/gps/waiver", gpsH.CreateWaiver)
		api.GET("/mobile/gps/waivers", gpsH.ListWaivers)
		api.DELETE("/mobile/gps/waivers/:id", gpsH.DeleteWaiver)
		api.GET("/visits/:id/audit", auditH.GetByVisit)
		api.GET("/properties/:id/history", propertyHistoryH.GetByProperty)
		api.GET("/properties/:id/partners", mobilePartnerH.GetByProperty)
		api.GET("/templates", mobilePartnerH.ListTemplates)
	}

	return r, cleanup
}
