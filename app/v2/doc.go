/*
Package app provides an application framework that manages the live cycle of a running application.

An app could be an HTTP server (or any kind of server), a worker or a simple program. Independently of the kind of the app, it always exposes an admin port at 9000 by default which serves metrics, debug information and kubernetes probes. It also is capable of graceful shutdown on receiving signals of terminating by itself.

The recommended way to use the `app` package is to rely on the 'default app'. The 'default app' is a global app that can be accessed by public functions on the app package (the app package can be seen as the application)

There is a running example on `app/examples/servefiles`

# Basics

The first thing you should do bootstrap the application. This will initialize the configuration, reading it from the commandline or environment, initialize zerolog (based on the Config) and initialize the default global app.

	app.Bootstrap(version, &cfg)

After the bootstrap, the app will already be exposing the admin port and the readiness probe will be returning error, indicating that the application is not yet ready to receive requests. But the liveless probe will be returning success, indicating the app is alive.

Then you should start initializing all the program dependencies. Because the application is not yet ready, kubernetes will refrain from sending requests (that would fail at this point). Also we already have some metrics and the debug handlers.

During this phase, you will probably want to register your shutdown handlers.

	app.RegisterShutdownHandler(
		&app.ShutdownHandler{
			Name:     "http_server",
			Priority: app.ShutdownPriority(100),
			Handler:  httpServer.Shutdown,
			Policy:   app.ErrorPolicyAbort,
		},
	)

They are executed in order by priority. The Highest priority first (in case of the same priority, don't assume any order).

Finally you can run the application by calling RunAndWait:

	app.RunAndWait(func(ctx context.Context) error {...})

At this point the application will run until the given function returns or it receives an termination signal.

# Updating From Previous Version

Version 2 is a major overhaul over version 1. One of the main breaking changes is how the Config struct is used. Now, all app related configuration is inside an App field and new configuration were added. Now the Config struct is expected to be embedded in your application's configuration:

	type config struct {
		// App is the app specific configuration
		app.Config

		// Programs can have any configuration the want.
		HTTP struct {
			Port string `default:"8000"`
		}
		Dir string `default:"."`
	}

All the initialization now occurs on the Bootstrap functions and you need to initialize logs manually anymore.

RunAndWait was also changed in a couple of ways. Now it does not return anything and will panic if called incorrectly. The error returned by the MainLoopFunc is handled my logging it and triggering a graceful shutdown. The MainLoopFunc was changed to receve a context. The context will be canceled when Shutdown is called.

# Using Probes

A Probe is a boolean that indicates if something is OK or not. There are two groups of probes in an app: The Healthiness an Readiness groups. Kubernetes checks on there two probes to decide what to do to the pod, like, from stop sending requests to just kill the pod, sending a signal the app will capture and start a graceful shutdown.

If a single probe of a group is not OK, than the whole group is not OK. In this event, the HTTP handler returns the name of all the probes that are not OK for the given group.

	mux.HandleFunc("/healthy", func(w http.ResponseWriter, _ *http.Request) {
		isHealthy, cause := app.Healthy.CheckProbes()
		if isHealthy {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(cause))
		}
	})
	mux.HandleFunc("/ready", func(w http.ResponseWriter, _ *http.Request) {
		isReady, cause := app.Ready.CheckProbes()
		if isReady {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(cause))
		}
	})

If the application is unhealthy kubernetes will send a signal that will trigger the graceful shutdown. All registered shutdown handlers will be executed ordered by priority (highest first) and the pod will be restarted. Only set an application as unhealthy if it reached an unworkable state and should be restarted. We have an example of this on `gokitmiddlewares/stalemiddleware/`. This is a middleware that was developed to be used in workers. It checks if the endpoint is being called (messages are being fetched and processed) and if not, it assumes there could be a problem with the queue and sets the application to unready, causing the application to restart. This mitigated a problem we had with kafka when a change of brokers made the worker stop receiving messages forever.

If the application is unready kubernetes will stop sending requests, but if the application becomes ready again, it will start receiving requests. This is used during initialization to signalize to kubernetes when the application is ready and can receive requests. If we can identify that the the application is degraded we can use this probe to temporary remove the application from the kubernetes service until it recovers.

A probe only exists as part of a group so the group provides a proper constructor for a probe. Probe's name must also be unique for the group but can be reused on different groups.

	readinessProbe, err := app.Ready.NewProbe("fkit/app", false)
	healthnessProbe, err := app.Healthy.NewProbe("fkit/app", true)

The probe is automatically added to the group and any change is automatically reflected on the group it belongs to and the HTTP probe endpoints.

The state of a probe can be altered at any time using SetOk and SetNotOk:

	readinessProbe.SetOk()
	readinessProbe.SetNotOk()
*/
package app
