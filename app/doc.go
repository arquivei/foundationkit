/*
Package app provides an application framework that manages the live cycle of a running application.

An app could be an HTTP server (or any kind of server), a worker or a simple program. Independently of the kind of the app, it always exposes an admin port at 9000 by default which serves metrics, debug information and kubernetes probes. It also is capable of graceful shutdown on receiving signals of terminating by itself.

The recommended way to use the `app` package is to rely on the 'default app'. The 'default app' is a global app that can be accessed by public functions on the app package (the app package can be seen as the application)

There is a running example on `app/examples/servefiles`

Basics

The first thing you should do is setup a configuration and logger. The app uses the foundation's kit log package and expects a zerolog's logger in the context.

	app.SetupConfig(&config)
	ctx := log.SetupLoggerWithContext(context.Background(), config.Log, version)
	app.NewDefaultApp(ctx)

At this point, the app will already be exposing the admin port and the readiness probe will be returning error, indicating that the application is not yet ready to receive requests.

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

	app.RunAndWait(func() error {
		return httpServer.ListenAndServe()
	})

At this point the application will run until the given function returns or it receives an termination signal.

Updating From Previous Version

On the previous version,the NewDefaultApp received the main loop:

   func NewDefaultApp(ctx context.Context, mainLoop MainLoopFunc) (err error)

This was a problem because the main loop normally depends on various resources that must be created before the main loop can be called. But the creation of this resourced involves registering shutdown handlers, that requires an already created app.

This cycle forced the application to rely on lazy initialization of the resources. Lazy initialization is not a bad thing but in this particular case this means that when we call RunAndWait and the readiness probe is set return success, the application is still initializing and could start receiving requests before it was really ready.

To break this cycle the main loop was moved from the NewDefaultApp and was placed on the RunAndWait function. So to update to this version you could only change these two functions calls. But, to really take advantage of this new way to start an app, you should refactor the code to remove the laziness part before the RunAndWait is called.

Using Probes

A Probe is a boolean that indicates if something is OK or not. There are two groups of probes in an app: The Healthiness an Readiness groups. Kubernetes checks on there two probes to decide what to do to the pod, like, from stop sending requests to just kill the pod, sending a signal the app will capture and start a graceful shutdown.

If a single probe of a group is not ok, than the whole group is not ok. In this event, the HTTP handler returns the name of all the probes that are not ok for the given group.

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
