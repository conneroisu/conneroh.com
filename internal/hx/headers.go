package hx

const (
	// Request Headers.

	// HdrBoosted indicates that the request was sent via an element using [hx-boost](https://htmx.org/docs/#boost).
	HdrBoosted = "HX-Boosted"
	// HdrCurrentURL indicates the current URL of the page.
	HdrCurrentURL = "HX-Current-Url"
	// HdrHistoryRestoreRequest indicates that the request was sent via an element using [hx-history-restore](https://htmx.org/docs/#history-restore).
	HdrHistoryRestoreRequest = "HX-History-Restore-Request"
	// HdrPrompt indicates that the request was sent via an element using [hx-prompt](https://htmx.org/docs/#prompt).
	HdrPrompt = "HX-Prompt"
	// HdrRequest is always set to "true".
	HdrRequest = "HX-Request"
	// HdrTarget is the id of the target element if it exists.
	HdrTarget = "HX-Target"
	// HdrTriggerName is the name of the triggered element if it exists.
	HdrTriggerName = "HX-Trigger-Name"
	// HdrTrigger is the id of the triggered element if it exists.
	HdrTrigger = "HX-Trigger"

	// Response Headers.

	// HdrLocation is the URL to redirect to without doing a full page reload.
	HdrLocation = "HX-Location"
	// HdrPushURL is the URL to pushState into client history to without doing a full page reload.
	HdrPushURL = "HX-Push-Url"
	// HdrRedirect is the URL to client-side redirect to.
	HdrRedirect = "HX-Redirect"
	// HdrRefresh is the URL to client-side refresh to.
	HdrRefresh = "HX-Refresh"
	// HdrReplaceURL is the URL to replaceState into client history to without doing a full page reload.
	HdrReplaceURL = "HX-Replace-Url"
	// HdrReswap allows you to specify how the response will be swapped.
	HdrReswap = "HX-Reswap"
	// HdrRetarget is a CSS selector that updates the target of the content update to a different element on the page.
	HdrRetarget = "HX-Retarget"
	// HdrReselect is a CSS selector that allows you to choose which part of the response is used to be swapped in.
	HdrReselect = "HX-Reselect"
	// HdrTriggerResponse allows you to trigger client-side events.
	HdrTriggerResponse = "HX-Trigger"
	// HdrTriggerAfterSettle allows you to trigger client-side events after the settle step.
	HdrTriggerAfterSettle = "HX-Trigger-After-Settle"
	// HdrTriggerAfterSwap allows you to trigger client-side events after the swap step.
	HdrTriggerAfterSwap = "HX-Trigger-After-Swap"
)
