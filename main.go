/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// +kubebuilder:docs-gen:collapse=Apache License

package main

import (
	batchv1 "example/api/v1"
	"example/controllers"
	"flag"
	"os"

	kbatchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

// +kubebuilder:docs-gen:collapse=Imports

/*
The first difference to notice is that kubebuilder has added the new API
group's package (`batchv1`) to our scheme.  This means that we can use those
objects in our controller.

We'll also need to add the kubernetes batch v1 (`kbatchv1`) scheme, since we're creating
and listing Jobs.
*/
var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)

	_ = kbatchv1.AddToScheme(scheme) // we've added this ourselves
	_ = batchv1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

/*
The other thing that's changed is that kubebuilder has added a block calling our
CronJob controller's `SetupWithManager` method.  Since we now use a `Scheme` as well,
we'll need to pass that to the reconciler ourselves.
*/

func main() {
	/*
	 */
	var metricsAddr string
	var enableLeaderElection bool
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(zap.New(func(o *zap.Options) {
		o.Development = true
	}))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}
	// +kubebuilder:docs-gen:collapse=old stuff

	if err = (&controllers.CronJobReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Captain"),
		Scheme: mgr.GetScheme(), // we've added this ourselves
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Captain")
		os.Exit(1)
	}

	/*
		We'll also set up webhooks for our type, which we'll talk about next.
		We just need to add them to the manager.  Since we might want to run
		the webhooks separately, or not run them when testing our controller
		locally, we'll put them behind an environment variable.

		We'll just make sure to set `ENABLE_WEBHOOKS=false` when we run locally.
	*/
	if os.Getenv("ENABLE_WEBHOOKS") != "false" {
		if err = (&batchv1.CronJob{}).SetupWebhookWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create webhook", "webhook", "Captain")
			os.Exit(1)
		}
	}
	// +kubebuilder:scaffold:builder

	/*
	 */

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
	// +kubebuilder:docs-gen:collapse=old stuff
}
