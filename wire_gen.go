// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package main

import (
	"github.com/devtron-labs/devtron/api/connector"
	"github.com/devtron-labs/devtron/api/restHandler"
	"github.com/devtron-labs/devtron/api/router"
	pubsub2 "github.com/devtron-labs/devtron/api/router/pubsub"
	"github.com/devtron-labs/devtron/api/sse"
	"github.com/devtron-labs/devtron/client/argocdServer"
	"github.com/devtron-labs/devtron/client/argocdServer/application"
	cluster3 "github.com/devtron-labs/devtron/client/argocdServer/cluster"
	repository2 "github.com/devtron-labs/devtron/client/argocdServer/repository"
	session2 "github.com/devtron-labs/devtron/client/argocdServer/session"
	"github.com/devtron-labs/devtron/client/events"
	"github.com/devtron-labs/devtron/client/gitSensor"
	"github.com/devtron-labs/devtron/client/grafana"
	client2 "github.com/devtron-labs/devtron/client/jira"
	"github.com/devtron-labs/devtron/client/lens"
	"github.com/devtron-labs/devtron/client/pubsub"
	"github.com/devtron-labs/devtron/internal/casbin"
	"github.com/devtron-labs/devtron/internal/sql/models"
	"github.com/devtron-labs/devtron/internal/sql/repository"
	"github.com/devtron-labs/devtron/internal/sql/repository/appWorkflow"
	"github.com/devtron-labs/devtron/internal/sql/repository/appstore"
	"github.com/devtron-labs/devtron/internal/sql/repository/appstore/chartGroup"
	"github.com/devtron-labs/devtron/internal/sql/repository/chartConfig"
	"github.com/devtron-labs/devtron/internal/sql/repository/cluster"
	"github.com/devtron-labs/devtron/internal/sql/repository/helper"
	"github.com/devtron-labs/devtron/internal/sql/repository/pipelineConfig"
	"github.com/devtron-labs/devtron/internal/sql/repository/security"
	"github.com/devtron-labs/devtron/internal/sql/repository/team"
	"github.com/devtron-labs/devtron/internal/util"
	"github.com/devtron-labs/devtron/internal/util/ArgoUtil"
	"github.com/devtron-labs/devtron/pkg/app"
	"github.com/devtron-labs/devtron/pkg/appClone"
	"github.com/devtron-labs/devtron/pkg/appClone/batch"
	appWorkflow2 "github.com/devtron-labs/devtron/pkg/appWorkflow"
	appstore2 "github.com/devtron-labs/devtron/pkg/appstore"
	cluster2 "github.com/devtron-labs/devtron/pkg/cluster"
	"github.com/devtron-labs/devtron/pkg/commonService"
	"github.com/devtron-labs/devtron/pkg/deploymentGroup"
	"github.com/devtron-labs/devtron/pkg/dex"
	"github.com/devtron-labs/devtron/pkg/event"
	"github.com/devtron-labs/devtron/pkg/git"
	jira2 "github.com/devtron-labs/devtron/pkg/jira"
	"github.com/devtron-labs/devtron/pkg/notifier"
	"github.com/devtron-labs/devtron/pkg/pipeline"
	"github.com/devtron-labs/devtron/pkg/projectManagementService/jira"
	security2 "github.com/devtron-labs/devtron/pkg/security"
	"github.com/devtron-labs/devtron/pkg/sso"
	team2 "github.com/devtron-labs/devtron/pkg/team"
	"github.com/devtron-labs/devtron/pkg/terminal"
	"github.com/devtron-labs/devtron/pkg/user"
	"github.com/devtron-labs/devtron/util/rbac"
	"github.com/devtron-labs/devtron/util/session"
)

import (
	_ "github.com/lib/pq"
)

// Injectors from Wire.go:

func InitializeApp() (*App, error) {
	sugaredLogger := util.NewSugardLogger()
	config, err := models.GetConfig()
	if err != nil {
		return nil, err
	}
	db, err := models.NewDbConnection(config, sugaredLogger)
	if err != nil {
		return nil, err
	}
	envConfigOverrideRepositoryImpl := chartConfig.NewEnvConfigOverrideRepository(db)
	pipelineOverrideRepositoryImpl := chartConfig.NewPipelineOverrideRepository(db)
	mergeUtil := &util.MergeUtil{
		Logger: sugaredLogger,
	}
	ciArtifactRepositoryImpl := repository.NewCiArtifactRepositoryImpl(db, sugaredLogger)
	gitConfig, err := util.GetGitConfig()
	if err != nil {
		return nil, err
	}
	gitServiceImpl := util.NewGitServiceImpl(gitConfig, sugaredLogger)
	gitClient, err := util.NewGitLabClient(gitConfig, sugaredLogger, gitServiceImpl)
	if err != nil {
		return nil, err
	}
	pipelineRepositoryImpl := pipelineConfig.NewPipelineRepositoryImpl(db, sugaredLogger)
	dbMigrationConfigRepositoryImpl := pipelineConfig.NewDbMigrationConfigRepositoryImpl(db, sugaredLogger)
	httpClient := util.NewHttpClient()
	eventClientConfig, err := client.GetEventClientConfig()
	if err != nil {
		return nil, err
	}
	pubSubClient, err := pubsub.NewPubSubClient(sugaredLogger)
	if err != nil {
		return nil, err
	}
	ciPipelineRepositoryImpl := pipelineConfig.NewCiPipelineRepositoryImpl(db, sugaredLogger)
	eventRESTClientImpl := client.NewEventRESTClientImpl(sugaredLogger, httpClient, eventClientConfig, pubSubClient, ciPipelineRepositoryImpl, pipelineRepositoryImpl)
	cdWorkflowRepositoryImpl := pipelineConfig.NewCdWorkflowRepositoryImpl(db, sugaredLogger)
	ciWorkflowRepositoryImpl := pipelineConfig.NewCiWorkflowRepositoryImpl(db, sugaredLogger)
	ciPipelineMaterialRepositoryImpl := pipelineConfig.NewCiPipelineMaterialRepositoryImpl(db, sugaredLogger)
	userRepositoryImpl := repository.NewUserRepositoryImpl(db)
	eventSimpleFactoryImpl := client.NewEventSimpleFactoryImpl(sugaredLogger, cdWorkflowRepositoryImpl, pipelineOverrideRepositoryImpl, ciWorkflowRepositoryImpl, ciPipelineMaterialRepositoryImpl, ciPipelineRepositoryImpl, pipelineRepositoryImpl, userRepositoryImpl)
	argocdServerConfig, err := argocdServer.GetConfig()
	if err != nil {
		return nil, err
	}
	settingsManager, err := session.SettingsManager(argocdServerConfig)
	if err != nil {
		return nil, err
	}
	argoCDSettings := session.CDSettingsManager(settingsManager)
	serviceClientImpl := application.NewApplicationClientImpl(argoCDSettings, sugaredLogger)
	acdAuthConfig, err := user.GetACDAuthConfig()
	if err != nil {
		return nil, err
	}
	userAuthRepositoryImpl := repository.NewUserAuthRepositoryImpl(db, sugaredLogger)
	dexConfig, err := dex.GetConfig()
	if err != nil {
		return nil, err
	}
	sessionManager := session.SessionManager(settingsManager, dexConfig)
	sessionServiceClientImpl := session2.NewSessionServiceClient(argoCDSettings)
	userAuthServiceImpl := user.NewUserAuthServiceImpl(userAuthRepositoryImpl, sessionManager, sessionServiceClientImpl, sugaredLogger, userRepositoryImpl)
	tokenCache := user.NewTokenCache(sugaredLogger, acdAuthConfig, userAuthServiceImpl)
	enforcer := casbin.Create()
	enforcerImpl := rbac.NewEnforcerImpl(enforcer, sessionManager, sugaredLogger)
	teamRepositoryImpl := team.NewTeamRepositoryImpl(db)
	appRepositoryImpl := pipelineConfig.NewAppRepositoryImpl(db)
	environmentRepositoryImpl := cluster.NewEnvironmentRepositoryImpl(db)
	enforcerUtilImpl := rbac.NewEnforcerUtilImpl(sugaredLogger, teamRepositoryImpl, appRepositoryImpl, environmentRepositoryImpl, pipelineRepositoryImpl, ciPipelineRepositoryImpl)
	roleGroupRepositoryImpl := repository.NewRoleGroupRepositoryImpl(db, sugaredLogger)
	userServiceImpl := user.NewUserServiceImpl(userAuthRepositoryImpl, sessionManager, sessionServiceClientImpl, sugaredLogger, userRepositoryImpl, roleGroupRepositoryImpl)
	appListingRepositoryQueryBuilder := helper.NewAppListingRepositoryQueryBuilder(sugaredLogger)
	appListingRepositoryImpl := repository.NewAppListingRepositoryImpl(sugaredLogger, db, appListingRepositoryQueryBuilder)
	pipelineConfigRepositoryImpl := chartConfig.NewPipelineConfigRepository(db)
	configMapRepositoryImpl := chartConfig.NewConfigMapRepositoryImpl(sugaredLogger, db)
	appLevelMetricsRepositoryImpl := repository.NewAppLevelMetricsRepositoryImpl(db, sugaredLogger)
	envLevelAppMetricsRepositoryImpl := repository.NewEnvLevelAppMetricsRepositoryImpl(db, sugaredLogger)
	chartRepositoryImpl := chartConfig.NewChartRepository(db)
	commonServiceImpl := commonService.NewCommonServiceImpl(sugaredLogger, chartRepositoryImpl, envConfigOverrideRepositoryImpl)
	imageScanDeployInfoRepositoryImpl := security.NewImageScanDeployInfoRepositoryImpl(db, sugaredLogger)
	imageScanHistoryRepositoryImpl := security.NewImageScanHistoryRepositoryImpl(db, sugaredLogger)
	argoK8sClientImpl := argocdServer.NewArgoK8sClientImpl(sugaredLogger)
	appServiceImpl := app.NewAppService(envConfigOverrideRepositoryImpl, pipelineOverrideRepositoryImpl, mergeUtil, sugaredLogger, ciArtifactRepositoryImpl, gitClient, pipelineRepositoryImpl, dbMigrationConfigRepositoryImpl, eventRESTClientImpl, eventSimpleFactoryImpl, serviceClientImpl, tokenCache, acdAuthConfig, enforcerImpl, enforcerUtilImpl, userServiceImpl, appListingRepositoryImpl, appRepositoryImpl, environmentRepositoryImpl, pipelineConfigRepositoryImpl, configMapRepositoryImpl, appLevelMetricsRepositoryImpl, envLevelAppMetricsRepositoryImpl, chartRepositoryImpl, ciPipelineMaterialRepositoryImpl, cdWorkflowRepositoryImpl, commonServiceImpl, imageScanDeployInfoRepositoryImpl, imageScanHistoryRepositoryImpl, argoK8sClientImpl)
	validate, err := util.IntValidator()
	if err != nil {
		return nil, err
	}
	materialRepositoryImpl := pipelineConfig.NewMaterialRepositoryImpl(db)
	gitSensorConfig, err := gitSensor.GetGitSensorConfig()
	if err != nil {
		return nil, err
	}
	gitSensorClientImpl, err := gitSensor.NewGitSensorSession(gitSensorConfig, sugaredLogger)
	if err != nil {
		return nil, err
	}
	ciConfig, err := pipeline.GetCiConfig()
	if err != nil {
		return nil, err
	}
	appWorkflowRepositoryImpl := appWorkflow.NewAppWorkflowRepositoryImpl(sugaredLogger, db)
	dbPipelineOrchestratorImpl := pipeline.NewDbPipelineOrchestrator(appRepositoryImpl, sugaredLogger, materialRepositoryImpl, pipelineRepositoryImpl, ciPipelineRepositoryImpl, ciPipelineMaterialRepositoryImpl, gitSensorClientImpl, ciConfig, appWorkflowRepositoryImpl, environmentRepositoryImpl)
	dockerArtifactStoreRepositoryImpl := repository.NewDockerArtifactStoreRepositoryImpl(db)
	utilMergeUtil := util.MergeUtil{
		Logger: sugaredLogger,
	}
	propertiesConfigServiceImpl := pipeline.NewPropertiesConfigServiceImpl(sugaredLogger, envConfigOverrideRepositoryImpl, chartRepositoryImpl, utilMergeUtil, environmentRepositoryImpl, dbPipelineOrchestratorImpl, serviceClientImpl, envLevelAppMetricsRepositoryImpl, appLevelMetricsRepositoryImpl)
	ciTemplateRepositoryImpl := pipelineConfig.NewCiTemplateRepositoryImpl(db, sugaredLogger)
	ecrConfig, err := pipeline.GetEcrConfig()
	if err != nil {
		return nil, err
	}
	imageScanResultRepositoryImpl := security.NewImageScanResultRepositoryImpl(db, sugaredLogger)
	pipelineBuilderImpl := pipeline.NewPipelineBuilderImpl(sugaredLogger, dbPipelineOrchestratorImpl, dockerArtifactStoreRepositoryImpl, materialRepositoryImpl, appRepositoryImpl, pipelineRepositoryImpl, propertiesConfigServiceImpl, ciTemplateRepositoryImpl, ciPipelineRepositoryImpl, serviceClientImpl, chartRepositoryImpl, ciArtifactRepositoryImpl, ecrConfig, envConfigOverrideRepositoryImpl, environmentRepositoryImpl, pipelineConfigRepositoryImpl, utilMergeUtil, appWorkflowRepositoryImpl, ciConfig, cdWorkflowRepositoryImpl, appServiceImpl, imageScanResultRepositoryImpl, gitClient, argoK8sClientImpl)
	clusterRepositoryImpl := cluster.NewClusterRepositoryImpl(db, sugaredLogger)
	grafanaClientConfig, err := grafana.GetGrafanaClientConfig()
	if err != nil {
		return nil, err
	}
	grafanaClientImpl := grafana.NewGrafanaClientImpl(sugaredLogger, httpClient, grafanaClientConfig)
	installedAppRepositoryImpl := appstore.NewInstalledAppRepositoryImpl(sugaredLogger, db)
	clusterInstalledAppsRepositoryImpl := appstore.NewClusterInstalledAppsRepositoryImpl(db, sugaredLogger)
	clusterServiceImpl := cluster2.NewClusterServiceImpl(clusterRepositoryImpl, environmentRepositoryImpl, grafanaClientImpl, sugaredLogger, installedAppRepositoryImpl, clusterInstalledAppsRepositoryImpl)
	k8sUtil := util.NewK8sUtil(sugaredLogger)
	environmentServiceImpl := cluster2.NewEnvironmentServiceImpl(environmentRepositoryImpl, clusterServiceImpl, sugaredLogger, k8sUtil, propertiesConfigServiceImpl, grafanaClientImpl)
	teamServiceImpl := team2.NewTeamServiceImpl(sugaredLogger, teamRepositoryImpl, pipelineBuilderImpl, environmentServiceImpl, userServiceImpl)
	cdConfig, err := pipeline.GetCdConfig()
	if err != nil {
		return nil, err
	}
	cdWorkflowServiceImpl := pipeline.NewCdWorkflowServiceImpl(sugaredLogger, environmentRepositoryImpl, cdConfig, appServiceImpl)
	deploymentGroupRepositoryImpl := repository.NewDeploymentGroupRepositoryImpl(sugaredLogger, db)
	cvePolicyRepositoryImpl := security.NewPolicyRepositoryImpl(db)
	workflowDagExecutorImpl := pipeline.NewWorkflowDagExecutorImpl(sugaredLogger, pipelineRepositoryImpl, cdWorkflowRepositoryImpl, pubSubClient, appServiceImpl, cdWorkflowServiceImpl, cdConfig, ciArtifactRepositoryImpl, ciPipelineRepositoryImpl, materialRepositoryImpl, pipelineOverrideRepositoryImpl, userServiceImpl, deploymentGroupRepositoryImpl, environmentRepositoryImpl, enforcerImpl, enforcerUtilImpl, tokenCache, acdAuthConfig, eventSimpleFactoryImpl, eventRESTClientImpl, cvePolicyRepositoryImpl, imageScanResultRepositoryImpl)
	deploymentGroupAppRepositoryImpl := repository.NewDeploymentGroupAppRepositoryImpl(sugaredLogger, db)
	deploymentGroupServiceImpl := deploymentGroup.NewDeploymentGroupServiceImpl(appRepositoryImpl, sugaredLogger, pipelineRepositoryImpl, ciPipelineRepositoryImpl, deploymentGroupRepositoryImpl, environmentRepositoryImpl, deploymentGroupAppRepositoryImpl, ciArtifactRepositoryImpl, appWorkflowRepositoryImpl, workflowDagExecutorImpl)
	pipelineTriggerRestHandlerImpl := restHandler.NewPipelineRestHandler(appServiceImpl, userServiceImpl, validate, enforcerImpl, teamServiceImpl, sugaredLogger, enforcerUtilImpl, workflowDagExecutorImpl, deploymentGroupServiceImpl)
	sseSSE := sse.NewSSE()
	helmRouterImpl := router.NewHelmRouter(pipelineTriggerRestHandlerImpl, sseSSE)
	chartWorkingDir := _wireChartWorkingDirValue
	chartTemplateServiceImpl := util.NewChartTemplateServiceImpl(sugaredLogger, chartWorkingDir, gitClient, gitServiceImpl, httpClient)
	chartRepoRepositoryImpl := chartConfig.NewChartRepoRepositoryImpl(db)
	refChartDir := _wireRefChartDirValue
	defaultChart := _wireDefaultChartValue
	repositoryServiceClientImpl := repository2.NewServiceClientImpl(argoCDSettings, sugaredLogger)
	chartRefRepositoryImpl := chartConfig.NewChartRefRepositoryImpl(db)
	chartServiceImpl := pipeline.NewChartServiceImpl(chartRepositoryImpl, sugaredLogger, chartTemplateServiceImpl, chartRepoRepositoryImpl, appRepositoryImpl, refChartDir, defaultChart, utilMergeUtil, repositoryServiceClientImpl, gitConfig, chartRefRepositoryImpl, envConfigOverrideRepositoryImpl, pipelineConfigRepositoryImpl, configMapRepositoryImpl, environmentRepositoryImpl, pipelineRepositoryImpl, appLevelMetricsRepositoryImpl, httpClient)
	dbMigrationServiceImpl := pipeline.NewDbMogrationService(sugaredLogger, dbMigrationConfigRepositoryImpl)
	workflowServiceImpl := pipeline.NewWorkflowServiceImpl(sugaredLogger, ciConfig)
	ciServiceImpl := pipeline.NewCiServiceImpl(sugaredLogger, workflowServiceImpl, ciPipelineMaterialRepositoryImpl, ciWorkflowRepositoryImpl, ciConfig, eventRESTClientImpl, eventSimpleFactoryImpl, mergeUtil, ciPipelineRepositoryImpl)
	ciLogServiceImpl := pipeline.NewCiLogServiceImpl(sugaredLogger, ciServiceImpl, ciConfig)
	ciHandlerImpl := pipeline.NewCiHandlerImpl(sugaredLogger, ciServiceImpl, ciPipelineMaterialRepositoryImpl, gitSensorClientImpl, ciWorkflowRepositoryImpl, workflowServiceImpl, ciLogServiceImpl, ciConfig, ciArtifactRepositoryImpl, userServiceImpl, eventRESTClientImpl, eventSimpleFactoryImpl, ciPipelineRepositoryImpl, appListingRepositoryImpl)
	gitProviderRepositoryImpl := repository.NewGitProviderRepositoryImpl(db)
	gitRegistryConfigImpl := pipeline.NewGitRegistryConfigImpl(sugaredLogger, gitProviderRepositoryImpl, gitSensorClientImpl)
	dockerRegistryConfigImpl := pipeline.NewDockerRegistryConfigImpl(dockerArtifactStoreRepositoryImpl, sugaredLogger)
	cdHandlerImpl := pipeline.NewCdHandlerImpl(sugaredLogger, cdConfig, userServiceImpl, cdWorkflowRepositoryImpl, cdWorkflowServiceImpl, ciLogServiceImpl, ciArtifactRepositoryImpl, ciPipelineMaterialRepositoryImpl, pipelineRepositoryImpl, environmentRepositoryImpl, ciWorkflowRepositoryImpl)
	configMapServiceImpl := pipeline.NewConfigMapServiceImpl(chartRepositoryImpl, sugaredLogger, chartRepoRepositoryImpl, utilMergeUtil, pipelineConfigRepositoryImpl, configMapRepositoryImpl, envConfigOverrideRepositoryImpl, commonServiceImpl, appRepositoryImpl)
	appWorkflowServiceImpl := appWorkflow2.NewAppWorkflowServiceImpl(sugaredLogger, appWorkflowRepositoryImpl, dbPipelineOrchestratorImpl, ciPipelineRepositoryImpl, pipelineRepositoryImpl)
	appListingViewBuilderImpl := app.NewAppListingViewBuilderImpl(sugaredLogger)
	linkoutsRepositoryImpl := repository.NewLinkoutsRepositoryImpl(sugaredLogger, db)
	appListingServiceImpl := app.NewAppListingServiceImpl(sugaredLogger, appListingRepositoryImpl, serviceClientImpl, appRepositoryImpl, appListingViewBuilderImpl, pipelineRepositoryImpl, linkoutsRepositoryImpl, appLevelMetricsRepositoryImpl, envLevelAppMetricsRepositoryImpl, cdWorkflowRepositoryImpl, pipelineOverrideRepositoryImpl)
	appCloneServiceImpl := appClone.NewAppCloneServiceImpl(sugaredLogger, pipelineBuilderImpl, materialRepositoryImpl, chartServiceImpl, configMapServiceImpl, appWorkflowServiceImpl, appListingServiceImpl, propertiesConfigServiceImpl)
	imageScanObjectMetaRepositoryImpl := security.NewImageScanObjectMetaRepositoryImpl(db, sugaredLogger)
	cveStoreRepositoryImpl := security.NewCveStoreRepositoryImpl(db, sugaredLogger)
	policyServiceImpl := security2.NewPolicyServiceImpl(environmentServiceImpl, sugaredLogger, appRepositoryImpl, pipelineOverrideRepositoryImpl, cvePolicyRepositoryImpl, clusterServiceImpl, pipelineRepositoryImpl, imageScanResultRepositoryImpl, imageScanDeployInfoRepositoryImpl, imageScanObjectMetaRepositoryImpl, httpClient, ciArtifactRepositoryImpl, ciConfig, imageScanHistoryRepositoryImpl, cveStoreRepositoryImpl)
	pipelineConfigRestHandlerImpl := restHandler.NewPipelineRestHandlerImpl(pipelineBuilderImpl, sugaredLogger, chartServiceImpl, propertiesConfigServiceImpl, dbMigrationServiceImpl, serviceClientImpl, userServiceImpl, teamServiceImpl, enforcerImpl, ciHandlerImpl, validate, gitSensorClientImpl, ciPipelineRepositoryImpl, pipelineRepositoryImpl, enforcerUtilImpl, environmentServiceImpl, gitRegistryConfigImpl, dockerRegistryConfigImpl, cdHandlerImpl, appCloneServiceImpl, appWorkflowServiceImpl, materialRepositoryImpl, policyServiceImpl, imageScanResultRepositoryImpl)
	appWorkflowRestHandlerImpl := restHandler.NewAppWorkflowRestHandlerImpl(sugaredLogger, userServiceImpl, appWorkflowServiceImpl, teamServiceImpl, enforcerImpl, pipelineBuilderImpl, appRepositoryImpl, enforcerUtilImpl)
	pipelineConfigRouterImpl := router.NewPipelineRouterImpl(pipelineConfigRestHandlerImpl, appWorkflowRestHandlerImpl)
	dbConfigRepositoryImpl := repository.NewDbConfigRepositoryImpl(db, sugaredLogger)
	dbConfigServiceImpl := pipeline.NewDbConfigService(dbConfigRepositoryImpl, sugaredLogger)
	migrateDbRestHandlerImpl := restHandler.NewMigrateDbRestHandlerImpl(dockerRegistryConfigImpl, sugaredLogger, gitRegistryConfigImpl, dbConfigServiceImpl, userServiceImpl, validate, dbMigrationServiceImpl, enforcerImpl)
	migrateDbRouterImpl := router.NewMigrateDbRouterImpl(migrateDbRestHandlerImpl)
	clusterAccountsRepositoryImpl := cluster.NewClusterAccountsRepositoryImpl(db)
	clusterAccountsServiceImpl := cluster2.NewClusterAccountsServiceImpl(clusterAccountsRepositoryImpl, environmentRepositoryImpl, clusterServiceImpl, sugaredLogger)
	clusterAccountsRestHandlerImpl := restHandler.NewClusterAccountsRestHandlerImpl(clusterAccountsServiceImpl, sugaredLogger, userServiceImpl)
	clusterAccountsRouterImpl := router.NewClusterAccountsRouterImpl(clusterAccountsRestHandlerImpl)
	appListingRestHandlerImpl := restHandler.NewAppListingRestHandlerImpl(serviceClientImpl, appListingServiceImpl, teamServiceImpl, enforcerImpl, pipelineBuilderImpl, sugaredLogger, enforcerUtilImpl, deploymentGroupServiceImpl, userServiceImpl)
	appListingRouterImpl := router.NewAppListingRouterImpl(appListingRestHandlerImpl)
	environmentRestHandlerImpl := restHandler.NewEnvironmentRestHandlerImpl(environmentServiceImpl, sugaredLogger, userServiceImpl, validate, enforcerImpl, enforcerUtilImpl, userAuthServiceImpl)
	environmentRouterImpl := router.NewEnvironmentRouterImpl(environmentRestHandlerImpl)
	clusterServiceClientImpl := cluster3.NewServiceClientImpl(argoCDSettings, sugaredLogger)
	refChartProxyDir := _wireRefChartProxyDirValue
	appStoreApplicationVersionRepositoryImpl := appstore.NewAppStoreApplicationVersionRepositoryImpl(sugaredLogger, db)
	appStoreRepositoryImpl := appstore.NewAppStoreRepositoryImpl(sugaredLogger, db)
	appStoreVersionValuesRepositoryImpl := appstore.NewAppStoreVersionValuesRepositoryImpl(sugaredLogger, db)
	appStoreValuesServiceImpl := appstore2.NewAppStoreValuesServiceImpl(sugaredLogger, appStoreRepositoryImpl, appStoreApplicationVersionRepositoryImpl, installedAppRepositoryImpl, userServiceImpl, appStoreVersionValuesRepositoryImpl, utilMergeUtil)
	chartGroupDeploymentRepositoryImpl := chartGroup.NewChartGroupDeploymentRepositoryImpl(db, sugaredLogger)
	installedAppServiceImpl, err := appstore2.NewInstalledAppServiceImpl(chartRepositoryImpl, sugaredLogger, chartRepoRepositoryImpl, utilMergeUtil, pipelineConfigRepositoryImpl, configMapRepositoryImpl, installedAppRepositoryImpl, chartTemplateServiceImpl, refChartProxyDir, gitConfig, repositoryServiceClientImpl, appStoreApplicationVersionRepositoryImpl, environmentRepositoryImpl, teamRepositoryImpl, gitClient, appRepositoryImpl, serviceClientImpl, appStoreValuesServiceImpl, pubSubClient, tokenCache, chartGroupDeploymentRepositoryImpl, environmentServiceImpl, clusterInstalledAppsRepositoryImpl, argoK8sClientImpl)
	if err != nil {
		return nil, err
	}
	clusterRestHandlerImpl := restHandler.NewClusterRestHandlerImpl(clusterServiceImpl, sugaredLogger, clusterServiceClientImpl, environmentServiceImpl, clusterAccountsServiceImpl, userServiceImpl, validate, enforcerImpl, installedAppServiceImpl)
	clusterRouterImpl := router.NewClusterRouterImpl(clusterRestHandlerImpl)
	clusterHelmConfigRepositoryImpl := cluster.NewClusterHelmConfigRepositoryImpl(db)
	clusterHelmConfigServiceImpl := cluster2.NewClusterHelmConfigServiceImpl(clusterHelmConfigRepositoryImpl, clusterServiceImpl, sugaredLogger)
	clusterHelmConfigRestHandlerImpl := restHandler.NewClusterHelmConfigRestHandlerImpl(clusterHelmConfigServiceImpl, sugaredLogger, userServiceImpl)
	clusterHelmConfigRouterImpl := router.NewClusterHelmConfigRouterImpl(clusterHelmConfigRestHandlerImpl)
	gitWebhookRepositoryImpl := repository.NewGitWebhookRepositoryImpl(db)
	gitWebhookServiceImpl := git.NewGitWebhookServiceImpl(sugaredLogger, ciHandlerImpl, gitWebhookRepositoryImpl)
	gitWebhookRestHandlerImpl := restHandler.NewGitWebhookRestHandlerImpl(sugaredLogger, gitWebhookServiceImpl)
	webhookServiceImpl := pipeline.NewWebhookServiceImpl(ciArtifactRepositoryImpl, sugaredLogger, ciPipelineRepositoryImpl, appServiceImpl, eventRESTClientImpl, eventSimpleFactoryImpl, ciWorkflowRepositoryImpl, workflowDagExecutorImpl, ciHandlerImpl)
	ciEventHandlerImpl := pubsub2.NewCiEventHandlerImpl(sugaredLogger, pubSubClient, webhookServiceImpl)
	externalCiRestHandlerImpl := restHandler.NewExternalCiRestHandlerImpl(sugaredLogger, webhookServiceImpl, ciEventHandlerImpl)
	natsPublishClientImpl := pubsub.NewNatsPublishClientImpl(sugaredLogger, pubSubClient)
	pubSubClientRestHandlerImpl := restHandler.NewPubSubClientRestHandlerImpl(natsPublishClientImpl, sugaredLogger, cdConfig)
	webhookRouterImpl := router.NewWebhookRouterImpl(gitWebhookRestHandlerImpl, pipelineConfigRestHandlerImpl, externalCiRestHandlerImpl, pubSubClientRestHandlerImpl)
	ssoLoginRepositoryImpl := repository.NewSSOLoginRepositoryImpl(db)
	ssoLoginServiceImpl := sso.NewSSOLoginServiceImpl(userAuthRepositoryImpl, sessionManager, sessionServiceClientImpl, sugaredLogger, userRepositoryImpl, roleGroupRepositoryImpl, ssoLoginRepositoryImpl, k8sUtil, clusterServiceImpl, environmentServiceImpl, acdAuthConfig)
	userAuthHandlerImpl := restHandler.NewUserAuthHandlerImpl(userAuthServiceImpl, validate, sugaredLogger, enforcerImpl, pubSubClient, userServiceImpl, ssoLoginServiceImpl)
	userAuthRouterImpl := router.NewUserAuthRouterImpl(sugaredLogger, userAuthHandlerImpl, argocdServerConfig, argoCDSettings, userServiceImpl)
	pumpImpl := connector.NewPumpImpl(sugaredLogger)
	terminalSessionHandlerImpl := terminal.NewTerminalSessionHandlerImpl(environmentServiceImpl, sugaredLogger)
	applicationRestHandlerImpl := restHandler.NewApplicationRestHandlerImpl(serviceClientImpl, pumpImpl, enforcerImpl, teamServiceImpl, environmentServiceImpl, sugaredLogger, enforcerUtilImpl, terminalSessionHandlerImpl)
	applicationRouterImpl := router.NewApplicationRouterImpl(applicationRestHandlerImpl, sugaredLogger)
	argoConfig, err := ArgoUtil.GetArgoConfig()
	if err != nil {
		return nil, err
	}
	argoSession, err := ArgoUtil.NewArgoSession(argoConfig, sugaredLogger)
	if err != nil {
		return nil, err
	}
	resourceServiceImpl := ArgoUtil.NewResourceServiceImpl(argoSession)
	cdRestHandlerImpl := restHandler.NewCDRestHandlerImpl(sugaredLogger, resourceServiceImpl)
	cdRouterImpl := router.NewCDRouterImpl(sugaredLogger, cdRestHandlerImpl)
	jiraAccountRepositoryImpl := repository.NewJiraAccountRepositoryImpl(db)
	jiraClientImpl := client2.NewJiraClientImpl(sugaredLogger, httpClient)
	accountServiceImpl := jira.NewAccountServiceImpl(sugaredLogger, jiraAccountRepositoryImpl, jiraClientImpl)
	accountValidatorImpl := jira.NewAccountValidatorImpl(sugaredLogger, jiraClientImpl)
	projectManagementServiceImpl := jira2.NewProjectManagementServiceImpl(sugaredLogger, accountServiceImpl, jiraAccountRepositoryImpl, accountValidatorImpl)
	jiraRestHandlerImpl := restHandler.NewJiraRestHandlerImpl(projectManagementServiceImpl, sugaredLogger, userServiceImpl, validate)
	projectManagementRouterImpl := router.NewProjectManagementRouterImpl(jiraRestHandlerImpl)
	gitProviderRestHandlerImpl := restHandler.NewGitProviderRestHandlerImpl(dockerRegistryConfigImpl, sugaredLogger, gitRegistryConfigImpl, dbConfigServiceImpl, userServiceImpl, validate, enforcerImpl, teamServiceImpl)
	gitProviderRouterImpl := router.NewGitProviderRouterImpl(gitProviderRestHandlerImpl)
	dockerRegRestHandlerImpl := restHandler.NewDockerRegRestHandlerImpl(dockerRegistryConfigImpl, sugaredLogger, gitRegistryConfigImpl, dbConfigServiceImpl, userServiceImpl, validate, enforcerImpl, teamServiceImpl)
	dockerRegRouterImpl := router.NewDockerRegRouterImpl(dockerRegRestHandlerImpl)
	notificationSettingsRepositoryImpl := repository.NewNotificationSettingsRepositoryImpl(db)
	notificationConfigBuilderImpl := notifier.NewNotificationConfigBuilderImpl(sugaredLogger)
	slackNotificationRepositoryImpl := repository.NewSlackNotificationRepositoryImpl(db)
	sesNotificationRepositoryImpl := repository.NewSESNotificationRepositoryImpl(db)
	notificationConfigServiceImpl := notifier.NewNotificationConfigServiceImpl(sugaredLogger, notificationSettingsRepositoryImpl, notificationConfigBuilderImpl, ciPipelineRepositoryImpl, pipelineRepositoryImpl, slackNotificationRepositoryImpl, sesNotificationRepositoryImpl, teamRepositoryImpl, environmentRepositoryImpl, appRepositoryImpl, userRepositoryImpl)
	slackNotificationServiceImpl := notifier.NewSlackNotificationServiceImpl(sugaredLogger, slackNotificationRepositoryImpl, teamServiceImpl, userRepositoryImpl, notificationSettingsRepositoryImpl)
	sesNotificationServiceImpl := notifier.NewSESNotificationServiceImpl(sugaredLogger, sesNotificationRepositoryImpl, teamServiceImpl)
	notificationRestHandlerImpl := restHandler.NewNotificationRestHandlerImpl(dockerRegistryConfigImpl, sugaredLogger, gitRegistryConfigImpl, dbConfigServiceImpl, userServiceImpl, validate, notificationConfigServiceImpl, slackNotificationServiceImpl, sesNotificationServiceImpl, enforcerImpl, teamServiceImpl, environmentServiceImpl, pipelineBuilderImpl, enforcerUtilImpl)
	notificationRouterImpl := router.NewNotificationRouterImpl(notificationRestHandlerImpl)
	teamRestHandlerImpl := restHandler.NewTeamRestHandlerImpl(sugaredLogger, teamServiceImpl, dbConfigServiceImpl, userServiceImpl, enforcerImpl, validate, enforcerUtilImpl, userAuthServiceImpl)
	teamRouterImpl := router.NewTeamRouterImpl(teamRestHandlerImpl)
	gitWebhookHandlerImpl := pubsub2.NewGitWebhookHandler(sugaredLogger, pubSubClient, gitWebhookServiceImpl)
	workflowStatusUpdateHandlerImpl := pubsub2.NewWorkflowStatusUpdateHandlerImpl(sugaredLogger, pubSubClient, ciHandlerImpl, cdHandlerImpl, eventSimpleFactoryImpl, eventRESTClientImpl, cdWorkflowRepositoryImpl)
	applicationStatusUpdateHandlerImpl := pubsub2.NewApplicationStatusUpdateHandlerImpl(sugaredLogger, pubSubClient, appServiceImpl, workflowDagExecutorImpl)
	roleGroupServiceImpl := user.NewRoleGroupServiceImpl(userAuthRepositoryImpl, sessionManager, sessionServiceClientImpl, sugaredLogger, userRepositoryImpl, roleGroupRepositoryImpl)
	userRestHandlerImpl := restHandler.NewUserRestHandlerImpl(userServiceImpl, validate, sugaredLogger, enforcerImpl, pubSubClient, roleGroupServiceImpl)
	userRouterImpl := router.NewUserRouterImpl(userRestHandlerImpl, dexConfig, argocdServerConfig, argoCDSettings)
	eventRepositoryImpl := repository.NewEventRepositoryImpl(sugaredLogger, db)
	deploymentFailureHandlerImpl := app.NewDeploymentFailureHandlerImpl(sugaredLogger, appListingServiceImpl, eventRESTClientImpl, eventSimpleFactoryImpl)
	eventServiceImpl := event.NewEventServiceImpl(sugaredLogger, eventRepositoryImpl, deploymentFailureHandlerImpl)
	cronBasedEventReceiverImpl := pubsub2.NewCronBasedEventReceiverImpl(sugaredLogger, pubSubClient, eventServiceImpl)
	chartRefRestHandlerImpl := restHandler.NewChartRefRestHandlerImpl(chartServiceImpl, sugaredLogger)
	chartRefRouterImpl := router.NewChartRefRouterImpl(chartRefRestHandlerImpl)
	configMapRestHandlerImpl := restHandler.NewConfigMapRestHandlerImpl(pipelineBuilderImpl, sugaredLogger, chartServiceImpl, userServiceImpl, teamServiceImpl, enforcerImpl, pipelineRepositoryImpl, enforcerUtilImpl, configMapServiceImpl)
	configMapRouterImpl := router.NewConfigMapRouterImpl(configMapRestHandlerImpl)
	versionServiceImpl := argocdServer.NewVersionServiceImpl(argoCDSettings, sugaredLogger)
	appStoreServiceImpl := appstore2.NewAppStoreServiceImpl(sugaredLogger, appStoreRepositoryImpl, appStoreApplicationVersionRepositoryImpl, installedAppRepositoryImpl, userServiceImpl, chartRepoRepositoryImpl, k8sUtil, clusterServiceImpl, environmentServiceImpl, versionServiceImpl, acdAuthConfig)
	appStoreRestHandlerImpl := restHandler.NewAppStoreRestHandlerImpl(sugaredLogger, userServiceImpl, appStoreServiceImpl, serviceClientImpl, teamServiceImpl, enforcerImpl, enforcerUtilImpl, validate, httpClient)
	installedAppRestHandlerImpl := restHandler.NewInstalledAppRestHandlerImpl(pipelineBuilderImpl, sugaredLogger, chartServiceImpl, userServiceImpl, teamServiceImpl, enforcerImpl, pipelineRepositoryImpl, enforcerUtilImpl, configMapServiceImpl, installedAppServiceImpl, validate)
	appStoreValuesRestHandlerImpl := restHandler.NewAppStoreValuesRestHandlerImpl(pipelineBuilderImpl, sugaredLogger, chartServiceImpl, userServiceImpl, teamServiceImpl, enforcerImpl, pipelineRepositoryImpl, enforcerUtilImpl, configMapServiceImpl, installedAppServiceImpl, appStoreValuesServiceImpl)
	appStoreRouterImpl := router.NewAppStoreRouterImpl(appStoreRestHandlerImpl, installedAppRestHandlerImpl, appStoreValuesRestHandlerImpl)
	lensConfig, err := lens.GetLensConfig()
	if err != nil {
		return nil, err
	}
	lensClientImpl, err := lens.NewLensClientImpl(lensConfig, sugaredLogger)
	if err != nil {
		return nil, err
	}
	releaseDataServiceImpl := app.NewReleaseDataServiceImpl(pipelineOverrideRepositoryImpl, sugaredLogger, ciPipelineMaterialRepositoryImpl, eventRESTClientImpl, lensClientImpl)
	releaseMetricsRestHandlerImpl := restHandler.NewReleaseMetricsRestHandlerImpl(sugaredLogger, enforcerImpl, releaseDataServiceImpl, userServiceImpl, teamServiceImpl, pipelineRepositoryImpl, enforcerUtilImpl)
	releaseMetricsRouterImpl := router.NewReleaseMetricsRouterImpl(sugaredLogger, releaseMetricsRestHandlerImpl)
	deploymentGroupRestHandlerImpl := restHandler.NewDeploymentGroupRestHandlerImpl(deploymentGroupServiceImpl, sugaredLogger, validate, enforcerImpl, teamServiceImpl, userServiceImpl, enforcerUtilImpl)
	deploymentGroupRouterImpl := router.NewDeploymentGroupRouterImpl(deploymentGroupRestHandlerImpl)
	buildActionImpl := batch.NewBuildActionImpl(pipelineBuilderImpl, sugaredLogger, appRepositoryImpl, appWorkflowRepositoryImpl, ciPipelineRepositoryImpl, materialRepositoryImpl)
	dataHolderActionImpl := batch.NewDataHolderActionImpl(appRepositoryImpl, configMapServiceImpl, environmentServiceImpl, sugaredLogger)
	deploymentTemplateActionImpl := batch.NewDeploymentTemplateActionImpl(sugaredLogger, appRepositoryImpl, chartServiceImpl)
	deploymentActionImpl := batch.NewDeploymentActionImpl(pipelineBuilderImpl, sugaredLogger, appRepositoryImpl, environmentServiceImpl, appWorkflowRepositoryImpl, ciPipelineRepositoryImpl, pipelineRepositoryImpl, dataHolderActionImpl, deploymentTemplateActionImpl)
	workflowActionImpl := batch.NewWorkflowActionImpl(sugaredLogger, appRepositoryImpl, appWorkflowServiceImpl, buildActionImpl, deploymentActionImpl)
	batchOperationRestHandlerImpl := restHandler.NewBatchOperationRestHandlerImpl(userServiceImpl, enforcerImpl, workflowActionImpl, teamServiceImpl, sugaredLogger)
	batchOperationRouterImpl := router.NewBatchOperationRouterImpl(batchOperationRestHandlerImpl, sugaredLogger)
	chartGroupEntriesRepositoryImpl := chartGroup.NewChartGroupEntriesRepositoryImpl(db, sugaredLogger)
	chartGroupReposotoryImpl := chartGroup.NewChartGroupReposotoryImpl(db, sugaredLogger)
	chartGroupServiceImpl := appstore2.NewChartGroupServiceImpl(chartGroupEntriesRepositoryImpl, chartGroupReposotoryImpl, sugaredLogger, chartGroupDeploymentRepositoryImpl, installedAppRepositoryImpl, appStoreVersionValuesRepositoryImpl)
	chartGroupRestHandlerImpl := restHandler.NewChartGroupRestHandlerImpl(chartGroupServiceImpl, sugaredLogger, userServiceImpl, enforcerImpl, enforcerUtilImpl, validate)
	chartGroupRouterImpl := router.NewChartGroupRouterImpl(chartGroupRestHandlerImpl)
	testSuitRestHandlerImpl := restHandler.NewTestSuitRestHandlerImpl(sugaredLogger, userServiceImpl, validate, enforcerImpl, enforcerUtilImpl, eventClientConfig, httpClient)
	testSuitRouterImpl := router.NewTestSuitRouterImpl(testSuitRestHandlerImpl)
	imageScanServiceImpl := security2.NewImageScanServiceImpl(sugaredLogger, imageScanHistoryRepositoryImpl, imageScanResultRepositoryImpl, imageScanObjectMetaRepositoryImpl, cveStoreRepositoryImpl, imageScanDeployInfoRepositoryImpl, userServiceImpl, teamRepositoryImpl, appRepositoryImpl, environmentServiceImpl, ciArtifactRepositoryImpl, policyServiceImpl, pipelineRepositoryImpl, installedAppRepositoryImpl, ciPipelineRepositoryImpl)
	imageScanRestHandlerImpl := restHandler.NewImageScanRestHandlerImpl(sugaredLogger, imageScanServiceImpl, userServiceImpl, enforcerImpl, enforcerUtilImpl, environmentServiceImpl)
	imageScanRouterImpl := router.NewImageScanRouterImpl(imageScanRestHandlerImpl)
	policyRestHandlerImpl := restHandler.NewPolicyRestHandlerImpl(sugaredLogger, policyServiceImpl, userServiceImpl, userAuthServiceImpl, enforcerImpl, enforcerUtilImpl, environmentServiceImpl)
	policyRouterImpl := router.NewPolicyRouterImpl(policyRestHandlerImpl)
	muxRouter := router.NewMuxRouter(sugaredLogger, helmRouterImpl, pipelineConfigRouterImpl, migrateDbRouterImpl, clusterAccountsRouterImpl, appListingRouterImpl, environmentRouterImpl, clusterRouterImpl, clusterHelmConfigRouterImpl, webhookRouterImpl, userAuthRouterImpl, applicationRouterImpl, cdRouterImpl, projectManagementRouterImpl, gitProviderRouterImpl, dockerRegRouterImpl, notificationRouterImpl, teamRouterImpl, gitWebhookHandlerImpl, workflowStatusUpdateHandlerImpl, applicationStatusUpdateHandlerImpl, ciEventHandlerImpl, pubSubClient, userRouterImpl, cronBasedEventReceiverImpl, chartRefRouterImpl, configMapRouterImpl, appStoreRouterImpl, releaseMetricsRouterImpl, deploymentGroupRouterImpl, batchOperationRouterImpl, chartGroupRouterImpl, testSuitRouterImpl, imageScanRouterImpl, policyRouterImpl)
	mainApp := NewApp(muxRouter, sugaredLogger, sseSSE, sessionManager, versionServiceImpl, enforcer, db, pubSubClient)
	return mainApp, nil
}

var (
	_wireChartWorkingDirValue  = util.ChartWorkingDir("/tmp/charts/")
	_wireRefChartDirValue      = pipeline.RefChartDir("scripts/devtron-reference-helm-charts")
	_wireDefaultChartValue     = pipeline.DefaultChart("reference-app-rolling")
	_wireRefChartProxyDirValue = appstore2.RefChartProxyDir("scripts/devtron-reference-helm-charts")
)
