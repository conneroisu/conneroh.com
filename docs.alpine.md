##### advanced-async
# Async
Alpine is built to support asynchronous functions in most places it supports standard ones.
For example, let's say you have a simple function called `getLabel()` that you use as the input to an `x-text` directive:
```js
function getLabel() {
    return 'Hello World!'
}
```
```alpine
<span x-text="getLabel()"></span>
```
Because `getLabel` is synchronous, everything works as expected.
Now let's pretend that `getLabel` makes a network request to retrieve the label and can't return one instantaneously (asynchronous). By making `getLabel` an async function, you can call it from Alpine using JavaScript's `await` syntax.
```js
async function getLabel() {
    let response = await fetch('/api/label')
    return await response.text()
}
```
```alpine
<span x-text="await getLabel()"></span>
```
Additionally, if you prefer calling methods in Alpine without the trailing parenthesis, you can leave them out and Alpine will detect that the provided function is async and handle it accordingly. For example:
```alpine
<span x-text="getLabel"></span>
```
##### advanced-csp
# CSP (Content-Security Policy) Build
In order for Alpine to be able to execute plain strings from HTML attributes as JavaScript expressions, for example `x-on:click="console.log()"`, it needs to rely on utilities that violate the "unsafe-eval" Content Security Policy that some applications may enforce for security purposes.
> Under the hood, Alpine doesn't actually use eval() itself because it's slow and problematic. Instead it uses Function declarations, which are much better, but still violate "unsafe-eval".
In order to accommodate environments where this CSP is necessary, Alpine offer's an alternate build that doesn't violate "unsafe-eval", but has a more restrictive syntax.
Initialize it within index.js (bundle):
```js
import Alpine from '@alpinejs/csp'
window.Alpine = Alpine
Alpine.start()
```
## Basic Example
To provide a glimpse of how using the CSP build might feel, here is a copy-pastable HTML file with a working counter component using a common CSP setup:
```alpine
<html>
    <head>
        <meta http-equiv="Content-Security-Policy" content="default-src 'self'; script-src 'nonce-a23gbfz9e'">
        <script defer nonce="a23gbfz9e" src="https://cdn.jsdelivr.net/npm/@alpinejs/csp@3.x.x/dist/cdn.min.js"></script>
    </head>
    <body>
        <div x-data="counter">
            <button x-on:click="increment"></button>
            <span x-text="count"></span>
        </div>
        <script nonce="a23gbfz9e">
            document.addEventListener('alpine:init', () => {
                Alpine.data('counter', () => {
                    return {
                        count: 1,
                        increment() {
                            this.count++;
                        },
                    }
                })
            })
        </script>
    </body>
</html>
```
## API Restrictions
Since Alpine can no longer interpret strings as plain JavaScript, it has to parse and construct JavaScript functions from them manually.
Due to this limitation, you must use `Alpine.data` to register your `x-data` objects, and must reference properties and methods from it by key only.
For example, an inline component like this will not work.
```alpine
<div x-data="{ count: 1 }">
    <button @click="count++">Increment</button>
    <span x-text="count"></span>
</div>
```
However, breaking out the expressions into external APIs, the following is valid with the CSP build:
```alpine
<div x-data="counter">
    <button @click="increment">Increment</button>
    <span x-text="count"></span>
</div>
```
```js
Alpine.data('counter', () => ({
    count: 1,
    increment() {
        this.count++
    },
}))
```
The CSP build supports accessing nested properties (property accessors) using the dot notation.
```alpine
<div x-data="counter">
    <button @click="foo.increment">Increment</button>
    <span x-text="foo.count"></span>
</div>
```
```js
Alpine.data('counter', () => ({
    foo: {
        count: 1,
        increment() {
            this.count++
        },
    },
}))
```
##### advanced-extending
# Extending
Alpine has a very open codebase that allows for extension in a number of ways. In fact, every available directive and magic in Alpine itself uses these exact APIs. In theory you could rebuild all of Alpine's functionality using them yourself.
## Lifecycle concerns
Before we dive into each individual API, let's first talk about where in your codebase you should consume these APIs.
Because these APIs have an impact on how Alpine initializes the page, they must be registered AFTER Alpine is downloaded and available on the page, but BEFORE it has initialized the page itself.
There are two different techniques depending on if you are importing Alpine into a bundle, or including it directly via a `<script>` tag. Let's look at them both:
### Via a script tag
If you are including Alpine via a script tag, you will need to register any custom extension code inside an `alpine:init` event listener.
Here's an example:
```alpine
<html>
    <script src="/js/alpine.js" defer></script>
    <div x-data x-foo></div>
    <script>
        document.addEventListener('alpine:init', () => {
            Alpine.directive('foo', ...)
        })
    </script>
</html>
```
If you want to extract your extension code into an external file, you will need to make sure that file's `<script>` tag is located BEFORE Alpine's like so:
```alpine
<html>
    <script src="/js/foo.js" defer></script>
    <script src="/js/alpine.js" defer></script>
    <div x-data x-foo></div>
</html>
```
### Via an NPM module
If you imported Alpine into a bundle, you have to make sure you are registering any extension code IN BETWEEN when you import the `Alpine` global object, and when you initialize Alpine by calling `Alpine.start()`. For example:
```js
import Alpine from 'alpinejs'
Alpine.directive('foo', ...)
window.Alpine = Alpine
window.Alpine.start()
```
Now that we know where to use these extension APIs, let's look more closely at how to use each one:
## Custom directives
Alpine allows you to register your own custom directives using the `Alpine.directive()` API.
### Method Signature
```js
Alpine.directive('[name]', (el, { value, modifiers, expression }, { Alpine, effect, cleanup }) => {})
```
- name - The name of the directive. The name "foo" for example would be consumed as `x-foo`
- el - The DOM element the directive is added to
- value - If provided, the part of the directive after a colon. Ex: `'bar'` in `x-foo:bar`
- modifiers - An array of dot-separated trailing additions to the directive. Ex: `['baz', 'lob']` from `x-foo.baz.lob`
- expression - The attribute value portion of the directive. Ex: `law` from `x-foo="law"`
- Alpine - The Alpine global object
- effect - A function to create reactive effects that will auto-cleanup after this directive is removed from the DOM
- cleanup - A function you can pass bespoke callbacks to that will run when this directive is removed from the DOM
### Simple Example
Here's an example of a simple directive we're going to create called: `x-uppercase`:
```js
Alpine.directive('uppercase', el => {
    el.textContent = el.textContent.toUpperCase()
})
```
```alpine
<div x-data>
    <span x-uppercase>Hello World!</span>
</div>
```
### Evaluating expressions
When registering a custom directive, you may want to evaluate a user-supplied JavaScript expression:
For example, let's say you wanted to create a custom directive as a shortcut to `console.log()`. Something like:
```alpine
<div x-data="{ message: 'Hello World!' }">
    <div x-log="message"></div>
</div>
```
You need to retrieve the actual value of `message` by evaluating it as a JavaScript expression with the `x-data` scope.
Fortunately, Alpine exposes its system for evaluating JavaScript expressions with an `evaluate()` API. Here's an example:
```js
Alpine.directive('log', (el, { expression }, { evaluate }) => {
    // expression === 'message'
    console.log(
        evaluate(expression)
    )
})
```
Now, when Alpine initializes the `<div x-log...>`, it will retrieve the expression passed into the directive ("message" in this case), and evaluate it in the context of the current element's Alpine component scope.
### Introducing reactivity
Building on the `x-log` example from before, let's say we wanted `x-log` to log the value of `message` and also log it if the value changes.
Given the following template:
```alpine
<div x-data="{ message: 'Hello World!' }">
    <div x-log="message"></div>
    <button @click="message = 'yolo'">Change</button>
</div>
```
We want "Hello World!" to be logged initially, then we want "yolo" to be logged after pressing the `<button>`.
We can adjust the implementation of `x-log` and introduce two new APIs to achieve this: `evaluateLater()` and `effect()`:
```js
Alpine.directive('log', (el, { expression }, { evaluateLater, effect }) => {
    let getThingToLog = evaluateLater(expression)
    effect(() => {
        getThingToLog(thingToLog => {
            console.log(thingToLog)
        })
    })
})
```
Let's walk through the above code, line by line.
```js
let getThingToLog = evaluateLater(expression)
```
Here, instead of immediately evaluating `message` and retrieving the result, we will convert the string expression ("message") into an actual JavaScript function that we can run at any time. If you're going to evaluate a JavaScript expression more than once, it is highly recommended to first generate a JavaScript function and use that rather than calling `evaluate()` directly. The reason being that the process to interpret a plain string as a JavaScript function is expensive and should be avoided when unnecessary.
```js
effect(() => {
    ...
})
```
By passing in a callback to `effect()`, we are telling Alpine to run the callback immediately, then track any dependencies it uses (`x-data` properties like `message` in our case). Now as soon as one of the dependencies changes, this callback will be re-run. This gives us our "reactivity".
You may recognize this functionality from `x-effect`. It is the same mechanism under the hood.
You may also notice that `Alpine.effect()` exists and wonder why we're not using it here. The reason is that the `effect` function provided via the method parameter has special functionality that cleans itself up when the directive is removed from the page for any reason.
For example, if for some reason the element with `x-log` on it got removed from the page, by using `effect()` instead of `Alpine.effect()` when the `message` property is changed, the value will no longer be logged to the console.
```js
getThingToLog(thingToLog => {
    console.log(thingToLog)
})
```
Now we will call `getThingToLog`, which if you recall is the actual JavaScript function version of the string expression: "message".
You might expect `getThingToCall()` to return the result right away, but instead Alpine requires you to pass in a callback to receive the result.
The reason for this is to support async expressions like `await getMessage()`. By passing in a "receiver" callback instead of getting the result immediately, you are allowing your directive to work with async expressions as well.
### Cleaning Up
Let's say you needed to register an event listener from a custom directive. After that directive is removed from the page for any reason, you would want to remove the event listener as well.
Alpine makes this simple by providing you with a `cleanup` function when registering custom directives.
Here's an example:
```js
Alpine.directive('...', (el, {}, { cleanup }) => {
    let handler = () => {}
    window.addEventListener('click', handler)
    cleanup(() => {
        window.removeEventListener('click', handler)
    })
})
```
Now if the directive is removed from this element or the element is removed itself, the event listener will be removed as well.
### Custom order
By default, any new directive will run after the majority of the standard ones (with the exception of `x-teleport`). This is usually acceptable but some times you might need to run your custom directive before another specific one.
This can be achieved by chaining the `.before() function to `Alpine.directive()` and specifying which directive needs to run after your custom one.
```js
Alpine.directive('foo', (el, { value, modifiers, expression }) => {
    Alpine.addScopeToNode(el, {foo: 'bar'})
}).before('bind')
```
```alpine
<div x-data>
    <span x-foo x-bind:foo="foo"></span>
</div>
```
> Note, the directive name must be written without the `x-` prefix (or any other custom prefix you may use).
## Custom magics
Alpine allows you to register custom "magics" (properties or methods) using `Alpine.magic()`. Any magic you register will be available to all your application's Alpine code with the `$` prefix.
### Method Signature
```js
Alpine.magic('[name]', (el, { Alpine }) => {})
```
- name - The name of the magic. The name "foo" for example would be consumed as `$foo`
- el - The DOM element the magic was triggered from
- Alpine - The Alpine global object
### Magic Properties
Here's a basic example of a "$now" magic helper to easily get the current time from anywhere in Alpine:
```js
Alpine.magic('now', () => {
    return (new Date).toLocaleTimeString()
})
```
```alpine
<span x-text="$now"></span>
```
Now the `<span>` tag will contain the current time, resembling something like "12:00:00 PM".
As you can see `$now` behaves like a static property, but under the hood is actually a getter that evaluates every time the property is accessed.
Because of this, you can implement magic "functions" by returning a function from the getter.
### Magic Functions
For example, if we wanted to create a `$clipboard()` magic function that accepts a string to copy to clipboard, we could implement it like so:
```js
Alpine.magic('clipboard', () => {
    return subject => navigator.clipboard.writeText(subject)
})
```
```alpine
<button @click="$clipboard('hello world')">Copy "Hello World"</button>
```
Now that accessing `$clipboard` returns a function itself, we can immediately call it and pass it an argument like we see in the template with `$clipboard('hello world')`.
You can use the more brief syntax (a double arrow function) for returning a function from a function if you'd prefer:
```js
Alpine.magic('clipboard', () => subject => {
    navigator.clipboard.writeText(subject)
})
```
##### advanced-reactivity
# Reactivity
Alpine is "reactive" in the sense that when you change a piece of data, everything that depends on that data "reacts" automatically to that change.
Every bit of reactivity that takes place in Alpine, happens because of two very important reactive functions in Alpine's core: `Alpine.reactive()`, and `Alpine.effect()`.
> Alpine uses VueJS's reactivity engine under the hood to provide these functions.
Understanding these two functions will give you super powers as an Alpine developer, but also just as a web developer in general.
## Alpine.reactive()
Let's first look at `Alpine.reactive()`. This function accepts a JavaScript object as its parameter and returns a "reactive" version of that object. For example:
```js
let data = { count: 1 }
let reactiveData = Alpine.reactive(data)
```
Under the hood, when `Alpine.reactive` receives `data`, it wraps it inside a custom JavaScript proxy.
A proxy is a special kind of object in JavaScript that can intercept "get" and "set" calls to a JavaScript object.
At face value, `reactiveData` should behave exactly like `data`. For example:
```js
console.log(data.count) // 1
console.log(reactiveData.count) // 1
reactiveData.count = 2
console.log(data.count) // 2
console.log(reactiveData.count) // 2
```
What you see here is that because `reactiveData` is a thin wrapper around `data`, any attempts to get or set a property will behave exactly as if you had interacted with `data` directly.
The main difference here is that any time you modify or retrieve (get or set) a value from `reactiveData`, Alpine is aware of it and can execute any other logic that depends on this data.
`Alpine.reactive` is only the first half of the story. `Alpine.effect` is the other half, let's dig in.
## Alpine.effect()
`Alpine.effect` accepts a single callback function. As soon as `Alpine.effect` is called, it will run the provided function, but actively look for any interactions with reactive data. If it detects an interaction (a get or set from the aforementioned reactive proxy) it will keep track of it and make sure to re-run the callback if any of reactive data changes in the future. For example:
```js
let data = Alpine.reactive({ count: 1 })
Alpine.effect(() => {
    console.log(data.count)
})
```
When this code is first run, "1" will be logged to the console. Any time `data.count` changes, it's value will be logged to the console again.
This is the mechanism that unlocks all of the reactivity at the core of Alpine.
To connect the dots further, let's look at a simple "counter" component example without using Alpine syntax at all, only using `Alpine.reactive` and `Alpine.effect`:
```alpine
<button>Increment</button>
Count: <span></span>
```
```js
let button = document.querySelector('button')
let span = document.querySelector('span')
let data = Alpine.reactive({ count: 1 })
Alpine.effect(() => {
    span.textContent = data.count
})
button.addEventListener('click', () => {
    data.count = data.count + 1
})
```
<div x-data="{ count: 1 }" class="demo">
    <button @click="count++">Increment</button>
    <div>Count: <span x-text="count"></span></div>
</div>
As you can see, you can make any data reactive, and you can also wrap any functionality in `Alpine.effect`.
This combination unlocks an incredibly powerful programming paradigm for web development. Run wild and free.
##### directives-bind
# x-bind
`x-bind` allows you to set HTML attributes on elements based on the result of JavaScript expressions.
For example, here's a component where we will use `x-bind` to set the placeholder value of an input.
```alpine
<div x-data="{ placeholder: 'Type here...' }">
    <input type="text" x-bind:placeholder="placeholder">
</div>
```
## Shorthand syntax
If `x-bind:` is too verbose for your liking, you can use the shorthand: `:`. For example, here is the same input element as above, but refactored to use the shorthand syntax.
```alpine
<input type="text" :placeholder="placeholder">
```
> Despite not being included in the above snippet, `x-bind` cannot be used if no parent element has `x-data` defined. 
## Binding classes
`x-bind` is most often useful for setting specific classes on an element based on your Alpine state.
Here's a simple example of a simple dropdown toggle, but instead of using `x-show`, we'll use a "hidden" class to toggle an element.
```alpine
<div x-data="{ open: false }">
    <button x-on:click="open = ! open">Toggle Dropdown</button>
    <div :class="open ? '' : 'hidden'">
        Dropdown Contents...
    </div>
</div>
```
Now, when `open` is `false`, the "hidden" class will be added to the dropdown.
### Shorthand conditionals
In cases like these, if you prefer a less verbose syntax you can use JavaScript's short-circuit evaluation instead of standard conditionals:
```alpine
<div :class="show ? '' : 'hidden'">
<div :class="show || 'hidden'">
```
The inverse is also available to you. Suppose instead of `open`, we use a variable with the opposite value: `closed`.
```alpine
<div :class="closed ? 'hidden' : ''">
<div :class="closed && 'hidden'">
```
### Class object syntax
Alpine offers an additional syntax for toggling classes if you prefer. By passing a JavaScript object where the classes are the keys and booleans are the values, Alpine will know which classes to apply and which to remove. For example:
```alpine
<div :class="{ 'hidden': ! show }">
```
This technique offers a unique advantage to other methods. When using object-syntax, Alpine will NOT preserve original classes applied to an element's `class` attribute.
For example, if you wanted to apply the "hidden" class to an element before Alpine loads, AND use Alpine to toggle its existence you can only achieve that behavior using object-syntax:
```alpine
<div class="hidden" :class="{ 'hidden': ! show }">
```
In case that confused you, let's dig deeper into how Alpine handles `x-bind:class` differently than other attributes.
### Special behavior
`x-bind:class` behaves differently than other attributes under the hood.
Consider the following case.
```alpine
<div class="opacity-50" :class="hide && 'hidden'">
```
If "class" were any other attribute, the `:class` binding would overwrite any existing class attribute, causing `opacity-50` to be overwritten by either `hidden` or `''`.
However, Alpine treats `class` bindings differently. It's smart enough to preserve existing classes on an element.
For example, if `hide` is true, the above example will result in the following DOM element:
```alpine
<div class="opacity-50 hidden">
```
If `hide` is false, the DOM element will look like:
```alpine
<div class="opacity-50">
```
This behavior should be invisible and intuitive to most users, but it is worth mentioning explicitly for the inquiring developer or any special cases that might crop up.
## Binding styles
Similar to the special syntax for binding classes with JavaScript objects, Alpine also offers an object-based syntax for binding `style` attributes.
Just like the class objects, this syntax is entirely optional. Only use it if it affords you some advantage.
```alpine
<div :style="{ color: 'red', display: 'flex' }">
<div style="color: red; display: flex;" ...>
```
Conditional inline styling is possible using expressions just like with x-bind:class. Short circuit operators can be used here as well by using a styles object as the second operand.
```alpine
<div x-bind:style="true && { color: 'red' }">
<div style="color: red;">
```
One advantage of this approach is being able to mix it in with existing styles on an element:
```alpine
<div style="padding: 1rem;" :style="{ color: 'red', display: 'flex' }">
<div style="padding: 1rem; color: red; display: flex;" ...>
```
And like most expressions in Alpine, you can always use the result of a JavaScript expression as the reference:
```alpine
<div x-data="{ styles: { color: 'red', display: 'flex' }}">
    <div :style="styles">
</div>
<div ...>
    <div style="color: red; display: flex;" ...>
</div>
```
## Binding Alpine Directives Directly
`x-bind` allows you to bind an object of different directives and attributes to an element.
The object keys can be anything you would normally write as an attribute name in Alpine. This includes Alpine directives and modifiers, but also plain HTML attributes. The object values are either plain strings, or in the case of dynamic Alpine directives, callbacks to be evaluated by Alpine.
```alpine
<div x-data="dropdown">
    <button x-bind="trigger">Open Dropdown</button>
    <span x-bind="dialogue">Dropdown Contents</span>
</div>
<script>
    document.addEventListener('alpine:init', () => {
        Alpine.data('dropdown', () => ({
            open: false,
            trigger: {
                ['x-ref']: 'trigger',
                ['@click']() {
                    this.open = true
                },
            },
            dialogue: {
                ['x-show']() {
                    return this.open
                },
                ['@click.outside']() {
                    this.open = false
                },
            },
        }))
    })
</script>
```
There are a couple of caveats to this usage of `x-bind`:
> When the directive being "bound" or "applied" is `x-for`, you should return a normal expression string from the callback. For example: `['x-for']() { return 'item in items' }`
##### directives-cloak
# x-cloak
Sometimes, when you're using AlpineJS for a part of your template, there is a "blip" where you might see your uninitialized template after the page loads, but before Alpine loads.
`x-cloak` addresses this scenario by hiding the element it's attached to until Alpine is fully loaded on the page.
For `x-cloak` to work however, you must add the following CSS to the page.
```css
[x-cloak] { display: none !important; }
```
The following example will hide the `<span>` tag until its `x-show` is specifically set to true, preventing any "blip" of the hidden element onto screen as Alpine loads.
```alpine
<span x-cloak x-show="false">This will not 'blip' onto screen at any point</span>
```
`x-cloak` doesn't just work on elements hidden by `x-show` or `x-if`: it also ensures that elements containing data are hidden until the data is correctly set. The following example will hide the `<span>` tag until Alpine has set its text content to the `message` property.
```alpine
<span x-cloak x-text="message"></span>
```
When Alpine loads on the page, it removes all `x-cloak` property from the element, which also removes the `display: none;` applied by CSS, therefore showing the element.
## Alternative to global syntax
If you'd like to achieve this same behavior, but avoid having to include a global style, you can use the following cool, but admittedly odd trick:
```alpine
<template x-if="true">
    <span x-text="message"></span>
</template>
```
This will achieve the same goal as `x-cloak` by just leveraging the way `x-if` works.
Because `<template>` elements are "hidden" in browsers by default, you won't see the `<span>` until Alpine has had a chance to render the `x-if="true"` and show it.
Again, this solution is not for everyone, but it's worth mentioning for special cases.
##### directives-data
# x-data
Everything in Alpine starts with the `x-data` directive.
`x-data` defines a chunk of HTML as an Alpine component and provides the reactive data for that component to reference.
Here's an example of a contrived dropdown component:
```alpine
<div x-data="{ open: false }">
    <button @click="open = ! open">Toggle Content</button>
    <div x-show="open">
        Content...
    </div>
</div>
```
Don't worry about the other directives in this example (`@click` and `x-show`), we'll get to those in a bit. For now, let's focus on `x-data`.
## Scope
Properties defined in an `x-data` directive are available to all element children. Even ones inside other, nested `x-data` components.
For example:
```alpine
<div x-data="{ foo: 'bar' }">
    <span x-text="foo"></span>
    <div x-data="{ bar: 'baz' }">
        <span x-text="foo"></span>
        <div x-data="{ foo: 'bob' }">
            <span x-text="foo"></span>
        </div>
    </div>
</div>
```
## Methods
Because `x-data` is evaluated as a normal JavaScript object, in addition to state, you can store methods and even getters.
For example, let's extract the "Toggle Content" behavior into a method on  `x-data`.
```alpine
<div x-data="{ open: false, toggle() { this.open = ! this.open } }">
    <button @click="toggle()">Toggle Content</button>
    <div x-show="open">
        Content...
    </div>
</div>
```
Notice the added `toggle() { this.open = ! this.open }` method on `x-data`. This method can now be called from anywhere inside the component.
You'll also notice the usage of `this.` to access state on the object itself. This is because Alpine evaluates this data object like any standard JavaScript object with a `this` context.
If you prefer, you can leave the calling parenthesis off of the `toggle` method completely. For example:
```alpine
<button @click="toggle()">...</button>
<button @click="toggle">...</button>
```
## Getters
JavaScript [getters](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Functions/get) are handy when the sole purpose of a method is to return data based on other state.
Think of them like "computed properties" (although, they are not cached like Vue's computed properties).
Let's refactor our component to use a getter called `isOpen` instead of accessing `open` directly.
```alpine
<div x-data="{
    open: false,
    get isOpen() { return this.open },
    toggle() { this.open = ! this.open },
}">
    <button @click="toggle()">Toggle Content</button>
    <div x-show="isOpen">
        Content...
    </div>
</div>
```
Notice the "Content" now depends on the `isOpen` getter instead of the `open` property directly.
In this case there is no tangible benefit. But in some cases, getters are helpful for providing a more expressive syntax in your components.
## Data-less components
Occasionally, you want to create an Alpine component, but you don't need any data.
In these cases, you can always pass in an empty object.
```alpine
<div x-data="{}">
```
However, if you wish, you can also eliminate the attribute value entirely if it looks better to you.
```alpine
<div x-data>
```
## Single-element components
Sometimes you may only have a single element inside your Alpine component, like the following:
```alpine
<div x-data="{ open: true }">
    <button @click="open = false" x-show="open">Hide Me</button>
</div>
```
In these cases, you can declare `x-data` directly on that single element:
```alpine
<button x-data="{ open: true }" @click="open = false" x-show="open">
    Hide Me
</button>
```
## Re-usable Data
If you find yourself duplicating the contents of `x-data`, or you find the inline syntax verbose, you can extract the `x-data` object out to a dedicated component using `Alpine.data`.
Here's a quick example:
```alpine
<div x-data="dropdown">
    <button @click="toggle">Toggle Content</button>
    <div x-show="open">
        Content...
    </div>
</div>
<script>
    document.addEventListener('alpine:init', () => {
        Alpine.data('dropdown', () => ({
            open: false,
            toggle() {
                this.open = ! this.open
            },
        }))
    })
</script>
```
##### directives-effect
# x-effect
`x-effect` is a useful directive for re-evaluating an expression when one of its dependencies change. You can think of it as a watcher where you don't have to specify what property to watch, it will watch all properties used within it.
If this definition is confusing for you, that's ok. It's better explained through an example:
```alpine
<div x-data="{ label: 'Hello' }" x-effect="console.log(label)">
    <button @click="label += ' World!'">Change Message</button>
</div>
```
When this component is loaded, the `x-effect` expression will be run and "Hello" will be logged into the console.
Because Alpine knows about any property references contained within `x-effect`, when the button is clicked and `label` is changed, the effect will be re-triggered and "Hello World!" will be logged to the console.
##### directives-for
# x-for
Alpine's `x-for` directive allows you to create DOM elements by iterating through a list. Here's a simple example of using it to create a list of colors based on an array.
```alpine
<ul x-data="{ colors: ['Red', 'Orange', 'Yellow'] }">
    <template x-for="color in colors">
        <li x-text="color"></li>
    </template>
</ul>
```
<div class="demo">
    <ul x-data="{ colors: ['Red', 'Orange', 'Yellow'] }">
        <template x-for="color in colors">
            <li x-text="color"></li>
        </template>
    </ul>
</div>
You may also pass objects to `x-for`.
```alpine
<ul x-data="{ car: { make: 'Jeep', model: 'Grand Cherokee', color: 'Black' } }">
    <template x-for="(value, index) in car">
        <li>
            <span x-text="index"></span>: <span x-text="value"></span>
        </li>
    </template>
</ul>
```
<div class="demo">
    <ul x-data="{ car: { make: 'Jeep', model: 'Grand Cherokee', color: 'Black' } }">
        <template x-for="(value, index) in car">
            <li>
                <span x-text="index"></span>: <span x-text="value"></span>
            </li>
        </template>
    </ul>
</div>
There are two rules worth noting about `x-for`:
> `x-for` MUST be declared on a `<template>` element.
> That `<template>` element MUST contain only one root element
## Keys
It is important to specify unique keys for each `x-for` iteration if you are going to be re-ordering items. Without dynamic keys, Alpine may have a hard time keeping track of what re-orders and will cause odd side-effects.
```alpine
<ul x-data="{ colors: [
    { id: 1, label: 'Red' },
    { id: 2, label: 'Orange' },
    { id: 3, label: 'Yellow' },
]}">
    <template x-for="color in colors" :key="color.id">
        <li x-text="color.label"></li>
    </template>
</ul>
```
Now if the colors are added, removed, re-ordered, or their "id"s change, Alpine will preserve or destroy the iterated `<li>`elements accordingly.
## Accessing indexes
If you need to access the index of each item in the iteration, you can do so using the `([item], [index]) in [items]` syntax like so:
```alpine
<ul x-data="{ colors: ['Red', 'Orange', 'Yellow'] }">
    <template x-for="(color, index) in colors">
        <li>
            <span x-text="index + ': '"></span>
            <span x-text="color"></span>
        </li>
    </template>
</ul>
```
You can also access the index inside a dynamic `:key` expression.
```alpine
<template x-for="(color, index) in colors" :key="index">
```
## Iterating over a range
If you need to simply loop `n` number of times, rather than iterate through an array, Alpine offers a short syntax.
```alpine
<ul>
    <template x-for="i in 10">
        <li x-text="i"></li>
    </template>
</ul>
```
`i` in this case can be named anything you like.
> Despite not being included in the above snippet, `x-for` cannot be used if no parent element has `x-data` defined. 
## Contents of a `<template>`
As mentioned above, an `<template>` tag must contain only one root element.
For example, the following code will not work:
```alpine
<template x-for="color in colors">
    <span>The next color is </span><span x-text="color">
</template>
```
but this code will work:
```alpine
<template x-for="color in colors">
    <p>
        <span>The next color is </span><span x-text="color">
    </p>
</template>
```
##### directives-html
# x-html
`x-html` sets the "innerHTML" property of an element to the result of a given expression.
> ⚠️ Only use on trusted content and never on user-provided content. ⚠️
> Dynamically rendering HTML from third parties can easily lead to XSS vulnerabilities.
Here's a basic example of using `x-html` to display a user's username.
```alpine
<div x-data="{ username: '<strong>calebporzio</strong>' }">
    Username: <span x-html="username"></span>
</div>
```
<div class="demo">
    <div x-data="{ username: '<strong>calebporzio</strong>' }">
        Username: <span x-html="username"></span>
    </div>
</div>
Now the `<span>` tag's inner HTML will be set to "<strong>calebporzio</strong>".
##### directives-id
# x-id
`x-id` allows you to declare a new "scope" for any new IDs generated using `$id()`. It accepts an array of strings (ID names) and adds a suffix to each `$id('...')` generated within it that is unique to other IDs on the page.
`x-id` is meant to be used in conjunction with the `$id(...)` magic.
[Visit the $id documentation](/magics/id) for a better understanding of this feature.
Here's a brief example of this directive in use:
```alpine
<div x-id="['text-input']">
    <label :for="$id('text-input')">Username</label>
    
    <input type="text" :id="$id('text-input')">
    
</div>
<div x-id="['text-input']">
    <label :for="$id('text-input')">Username</label>
    
    <input type="text" :id="$id('text-input')">
    
</div>
```
> Despite not being included in the above snippet, `x-id` cannot be used if no parent element has `x-data` defined. 
##### directives-if
# x-if
`x-if` is used for toggling elements on the page, similarly to `x-show`, however it completely adds and removes the element it's applied to rather than just changing its CSS display property to "none".
Because of this difference in behavior, `x-if` should not be applied directly to the element, but instead to a `<template>` tag that encloses the element. This way, Alpine can keep a record of the element once it's removed from the page.
```alpine
<template x-if="open">
    <div>Contents...</div>
</template>
```
> Despite not being included in the above snippet, `x-if` cannot be used if no parent element has `x-data` defined.
## Caveats
Unlike `x-show`, `x-if`, does NOT support transitioning toggles with `x-transition`.
`<template>` tags can only contain one root element.
##### directives-ignore
# x-ignore
By default, Alpine will crawl and initialize the entire DOM tree of an element containing `x-init` or `x-data`.
If for some reason, you don't want Alpine to touch a specific section of your HTML, you can prevent it from doing so using `x-ignore`.
```alpine
<div x-data="{ label: 'From Alpine' }">
    <div x-ignore>
        <span x-text="label"></span>
    </div>
</div>
```
In the above example, the `<span>` tag will not contain "From Alpine" because we told Alpine to ignore the contents of the `div` completely.
##### directives-init
# x-init
The `x-init` directive allows you to hook into the initialization phase of any element in Alpine.
```alpine
<div x-init="console.log('I\'m being initialized!')"></div>
```
In the above example, "I\'m being initialized!" will be output in the console before it makes further DOM updates.
Consider another example where `x-init` is used to fetch some JSON and store it in `x-data` before the component is processed.
```alpine
<div
    x-data="{ posts: [] }"
    x-init="posts = await (await fetch('/posts')).json()"
>...</div>
```
## $nextTick
Sometimes, you want to wait until after Alpine has completely finished rendering to execute some code.
This would be something like `useEffect(..., [])` in react, or `mount` in Vue.
By using Alpine's internal `$nextTick` magic, you can make this happen.
```alpine
<div x-init="$nextTick(() => { ... })"></div>
```
## Standalone `x-init`
You can add `x-init` to any elements inside or outside an `x-data` HTML block. For example:
```alpine
<div x-data>
    <span x-init="console.log('I can initialize')"></span>
</div>
<span x-init="console.log('I can initialize too')"></span>
```
## Auto-evaluate init() method
If the `x-data` object of a component contains an `init()` method, it will be called automatically. For example:
```alpine
<div x-data="{
    init() {
        console.log('I am called automatically')
    }
}">
    ...
</div>
```
This is also the case for components that were registered using the `Alpine.data()` syntax.
```js
Alpine.data('dropdown', () => ({
    init() {
        console.log('I will get evaluated when initializing each "dropdown" component.')
    },
}))
```
If you have both an `x-data` object containing an `init()` method and an `x-init` directive, the `x-data` method will be called before the directive.
```alpine
<div
    x-data="{
        init() {
            console.log('I am called first')
        }
    }"
    x-init="console.log('I am called second')"
    >
    ...
</div>
```
##### directives-model
# x-model
`x-model` allows you to bind the value of an input element to Alpine data.
Here's a simple example of using `x-model` to bind the value of a text field to a piece of data in Alpine.
```alpine
<div x-data="{ message: '' }">
    <input type="text" x-model="message">
    <span x-text="message"></span>
</div>
```
<div class="demo">
    <div x-data="{ message: '' }">
        <input type="text" x-model="message" placeholder="Type message...">
        <div class="pt-4" x-text="message"></div>
    </div>
</div>
Now as the user types into the text field, the `message` will be reflected in the `<span>` tag.
`x-model` is two-way bound, meaning it both "sets" and "gets". In addition to changing data, if the data itself changes, the element will reflect the change.
We can use the same example as above but this time, we'll add a button to change the value of the `message` property.
```alpine
<div x-data="{ message: '' }">
    <input type="text" x-model="message">
    <button x-on:click="message = 'changed'">Change Message</button>
</div>
```
<div class="demo">
    <div x-data="{ message: '' }">
        <input type="text" x-model="message" placeholder="Type message...">
        <button x-on:click="message = 'changed'">Change Message</button>
    </div>
</div>
Now when the `<button>` is clicked, the input element's value will instantly be updated to "changed".
`x-model` works with the following input elements:
* `<input type="text">`
* `<textarea>`
* `<input type="checkbox">`
* `<input type="radio">`
* `<select>`
* `<input type="range">`
## Text inputs
```alpine
<input type="text" x-model="message">
<span x-text="message"></span>
```
<div class="demo">
    <div x-data="{ message: '' }">
        <input type="text" x-model="message" placeholder="Type message">
        <div class="pt-4" x-text="message"></div>
    </div>
</div>
> Despite not being included in the above snippet, `x-model` cannot be used if no parent element has `x-data` defined. 
## Textarea inputs
```alpine
<textarea x-model="message"></textarea>
<span x-text="message"></span>
```
<div class="demo">
    <div x-data="{ message: '' }">
        <textarea x-model="message" placeholder="Type message"></textarea>
        <div class="pt-4" x-text="message"></div>
    </div>
</div>
## Checkbox inputs
### Single checkbox with boolean
```alpine
<input type="checkbox" id="checkbox" x-model="show">
<label for="checkbox" x-text="show"></label>
```
<div class="demo">
    <div x-data="{ open: '' }">
        <input type="checkbox" id="checkbox" x-model="open">
        <label for="checkbox" x-text="open"></label>
    </div>
</div>
### Multiple checkboxes bound to array
```alpine
<input type="checkbox" value="red" x-model="colors">
<input type="checkbox" value="orange" x-model="colors">
<input type="checkbox" value="yellow" x-model="colors">
Colors: <span x-text="colors"></span>
```
<div class="demo">
    <div x-data="{ colors: [] }">
        <input type="checkbox" value="red" x-model="colors">
        <input type="checkbox" value="orange" x-model="colors">
        <input type="checkbox" value="yellow" x-model="colors">
        <div class="pt-4">Colors: <span x-text="colors"></span></div>
    </div>
</div>
## Radio inputs
```alpine
<input type="radio" value="yes" x-model="answer">
<input type="radio" value="no" x-model="answer">
Answer: <span x-text="answer"></span>
```
<div class="demo">
    <div x-data="{ answer: '' }">
        <input type="radio" value="yes" x-model="answer">
        <input type="radio" value="no" x-model="answer">
        <div class="pt-4">Answer: <span x-text="answer"></span></div>
    </div>
</div>
## Select inputs
### Single select
```alpine
<select x-model="color">
    <option>Red</option>
    <option>Orange</option>
    <option>Yellow</option>
</select>
Color: <span x-text="color"></span>
```
<div class="demo">
    <div x-data="{ color: '' }">
        <select x-model="color">
            <option>Red</option>
            <option>Orange</option>
            <option>Yellow</option>
        </select>
        <div class="pt-4">Color: <span x-text="color"></span></div>
    </div>
</div>
### Single select with placeholder
```alpine
<select x-model="color">
    <option value="" disabled>Select A Color</option>
    <option>Red</option>
    <option>Orange</option>
    <option>Yellow</option>
</select>
Color: <span x-text="color"></span>
```
<div class="demo">
    <div x-data="{ color: '' }">
        <select x-model="color">
            <option value="" disabled>Select A Color</option>
            <option>Red</option>
            <option>Orange</option>
            <option>Yellow</option>
        </select>
        <div class="pt-4">Color: <span x-text="color"></span></div>
    </div>
</div>
### Multiple select
```alpine
<select x-model="color" multiple>
    <option>Red</option>
    <option>Orange</option>
    <option>Yellow</option>
</select>
Colors: <span x-text="color"></span>
```
<div class="demo">
    <div x-data="{ color: '' }">
        <select x-model="color" multiple>
            <option>Red</option>
            <option>Orange</option>
            <option>Yellow</option>
        </select>
        <div class="pt-4">Color: <span x-text="color"></span></div>
    </div>
</div>
### Dynamically populated Select Options
```alpine
<select x-model="color">
    <template x-for="color in ['Red', 'Orange', 'Yellow']">
        <option x-text="color"></option>
    </template>
</select>
Color: <span x-text="color"></span>
```
<div class="demo">
    <div x-data="{ color: '' }">
        <select x-model="color">
            <template x-for="color in ['Red', 'Orange', 'Yellow']">
                <option x-text="color"></option>
            </template>
        </select>
        <div class="pt-4">Color: <span x-text="color"></span></div>
    </div>
</div>
## Range inputs
```alpine
<input type="range" x-model="range" min="0" max="1" step="0.1">
<span x-text="range"></span>
```
<div class="demo">
    <div x-data="{ range: 0.5 }">
        <input type="range" x-model="range" min="0" max="1" step="0.1">
        <div class="pt-4" x-text="range"></div>
    </div>
</div>
## Modifiers
### `.lazy`
On text inputs, by default, `x-model` updates the property on every keystroke. By adding the `.lazy` modifier, you can force an `x-model` input to only update the property when user focuses away from the input element.
This is handy for things like real-time form-validation where you might not want to show an input validation error until the user "tabs" away from a field.
```alpine
<input type="text" x-model.lazy="username">
<span x-show="username.length > 20">The username is too long.</span>
```
### `.number`
By default, any data stored in a property via `x-model` is stored as a string. To force Alpine to store the value as a JavaScript number, add the `.number` modifier.
```alpine
<input type="text" x-model.number="age">
<span x-text="typeof age"></span>
```
### `.boolean`
By default, any data stored in a property via `x-model` is stored as a string. To force Alpine to store the value as a JavaScript boolean, add the `.boolean` modifier. Both integers (1/0) and strings (true/false) are valid boolean values.
```alpine
<select x-model.boolean="isActive">
    <option value="true">Yes</option>
    <option value="false">No</option>
</select>
<span x-text="typeof isActive"></span>
```
### `.debounce`
By adding `.debounce` to `x-model`, you can easily debounce the updating of bound input.
This is useful for things like real-time search inputs that fetch new data from the server every time the search property changes.
```alpine
<input type="text" x-model.debounce="search">
```
The default debounce time is 250 milliseconds, you can easily customize this by adding a time modifier like so.
```alpine
<input type="text" x-model.debounce.500ms="search">
```
### `.throttle`
Similar to `.debounce` you can limit the property update triggered by `x-model` to only updating on a specified interval.
<input type="text" x-model.throttle="search">
The default throttle interval is 250 milliseconds, you can easily customize this by adding a time modifier like so.
```alpine
<input type="text" x-model.throttle.500ms="search">
```
### `.fill`
By default, if an input has a value attribute, it is ignored by Alpine and instead, the value of the input is set to the value of the property bound using `x-model`.
But if a bound property is empty, then you can use an input's value attribute to populate the property by adding the `.fill` modifier.
<div x-data="{ message: null }">
  <input type="text" x-model.fill="message" value="This is the default message.">
</div>
## Programmatic access
Alpine exposes under-the-hood utilities for getting and setting properties bound with `x-model`. This is useful for complex Alpine utilities that may want to override the default x-model behavior, or instances where you want to allow `x-model` on a non-input element.
You can access these utilities through a property called `_x_model` on the `x-model`ed element. `_x_model` has two methods to get and set the bound property:
* `el._x_model.get()` (returns the value of the bound property)
* `el._x_model.set()` (sets the value of the bound property)
```alpine
<div x-data="{ username: 'calebporzio' }">
    <div x-ref="div" x-model="username"></div>
    <button @click="$refs.div._x_model.set('phantomatrix')">
        Change username to: 'phantomatrix'
    </button>
    <span x-text="$refs.div._x_model.get()"></span>
</div>
```
<div class="demo">
    <div x-data="{ username: 'calebporzio' }">
        <div x-ref="div" x-model="username"></div>
        <button @click="$refs.div._x_model.set('phantomatrix')">
            Change username to: 'phantomatrix'
        </button>
        <span x-text="$refs.div._x_model.get()"></span>
    </div>
</div>
##### directives-modelable
# x-modelable
`x-modelable` allows you to expose any Alpine property as the target of the `x-model` directive.
Here's a simple example of using `x-modelable` to expose a variable for binding with `x-model`.
```alpine
<div x-data="{ number: 5 }">
    <div x-data="{ count: 0 }" x-modelable="count" x-model="number">
        <button @click="count++">Increment</button>
    </div>
    Number: <span x-text="number"></span>
</div>
```
<div class="demo">
    <div x-data="{ number: 5 }">
        <div x-data="{ count: 0 }" x-modelable="count" x-model="number">
            <button @click="count++">Increment</button>
        </div>
        Number: <span x-text="number"></span>
    </div>
</div>
As you can see the outer scope property "number" is now bound to the inner scope property "count".
Typically this feature would be used in conjunction with a backend templating framework like Laravel Blade. It's useful for abstracting away Alpine components into backend templates and exposing state to the outside through `x-model` as if it were a native input.
##### directives-on
# x-on
`x-on` allows you to easily run code on dispatched DOM events.
Here's an example of simple button that shows an alert when clicked.
```alpine
<button x-on:click="alert('Hello World!')">Say Hi</button>
```
> `x-on` can only listen for events with lower case names, as HTML attributes are case-insensitive. Writing `x-on:CLICK` will listen for an event named `click`. If you need to listen for a custom event with a camelCase name, you can use the [`.camel` helper](#camel) to work around this limitation. Alternatively, you can use [`x-bind`](/directives/bind#bind-directives) to attach an `x-on` directive to an element in javascript code (where case will be preserved).
## Shorthand syntax
If `x-on:` is too verbose for your tastes, you can use the shorthand syntax: `@`.
Here's the same component as above, but using the shorthand syntax instead:
```alpine
<button @click="alert('Hello World!')">Say Hi</button>
```
> Despite not being included in the above snippet, `x-on` cannot be used if no parent element has `x-data` defined.
## The event object
If you wish to access the native JavaScript event object from your expression, you can use Alpine's magic `$event` property.
```alpine
<button @click="alert($event.target.getAttribute('message'))" message="Hello World">Say Hi</button>
```
In addition, Alpine also passes the event object to any methods referenced without trailing parenthesis. For example:
```alpine
<button @click="handleClick">...</button>
<script>
    function handleClick(e) {
        // Now you can access the event object (e) directly
    }
</script>
```
## Keyboard events
Alpine makes it easy to listen for `keydown` and `keyup` events on specific keys.
Here's an example of listening for the `Enter` key inside an input element.
```alpine
<input type="text" @keyup.enter="alert('Submitted!')">
```
You can also chain these key modifiers to achieve more complex listeners.
Here's a listener that runs when the `Shift` key is held and `Enter` is pressed, but not when `Enter` is pressed alone.
```alpine
<input type="text" @keyup.shift.enter="alert('Submitted!')">
```
You can directly use any valid key names exposed via [`KeyboardEvent.key`](https://developer.mozilla.org/en-US/docs/Web/API/KeyboardEvent/key/Key_Values) as modifiers by converting them to kebab-case.
```alpine
<input type="text" @keyup.page-down="alert('Submitted!')">
```
For easy reference, here is a list of common keys you may want to listen for.
| Modifier                       | Keyboard Key                       |
| ------------------------------ | ---------------------------------- |
| `.shift`                       | Shift                              |
| `.enter`                       | Enter                              |
| `.space`                       | Space                              |
| `.ctrl`                        | Ctrl                               |
| `.cmd`                         | Cmd                                |
| `.meta`                        | Cmd on Mac, Windows key on Windows |
| `.alt`                         | Alt                                |
| `.up` `.down` `.left` `.right` | Up/Down/Left/Right arrows          |
| `.escape`                      | Escape                             |
| `.tab`                         | Tab                                |
| `.caps-lock`                   | Caps Lock                          |
| `.equal`                       | Equal, `=`                         |
| `.period`                      | Period, `.`                        |
| `.comma`                       | Comma, `,`                         |
| `.slash`                       | Forward Slash, `/`                 |
## Mouse events
Like the above Keyboard Events, Alpine allows the use of some key modifiers for handling `click` events.
| Modifier | Event Key |
| -------- | --------- |
| `.shift` | shiftKey  |
| `.ctrl`  | ctrlKey   |
| `.cmd`   | metaKey   |
| `.meta`  | metaKey   |
| `.alt`   | altKey    |
These work on `click`, `auxclick`, `context` and `dblclick` events, and even `mouseover`, `mousemove`, `mouseenter`, `mouseleave`, `mouseout`, `mouseup` and `mousedown`.
Here's an example of a button that changes behaviour when the `Shift` key is held down.
```alpine
<button type="button"
    @click="message = 'selected'"
    @click.shift="message = 'added to selection'">
    @mousemove.shift="message = 'add to selection'"
    @mouseout="message = 'select'"
    x-text="message"></button>
```
<div class="demo">
    <div x-data="{ message: '' }">
        <button type="button"
            @click="message = 'selected'"
            @click.shift="message = 'added to selection'"
            @mousemove.shift="message = 'add to selection'"
            @mouseout="message = 'select'"
            x-text="message"></button>
    </div>
</div>
> Note: Normal click events with some modifiers (like `ctrl`) will automatically become `contextmenu` events in most browsers. Similarly, `right-click` events will trigger a `contextmenu` event, but will also trigger an `auxclick` event if the `contextmenu` event is prevented.
## Custom events
Alpine event listeners are a wrapper for native DOM event listeners. Therefore, they can listen for ANY DOM event, including custom events.
Here's an example of a component that dispatches a custom DOM event and listens for it as well.
```alpine
<div x-data @foo="alert('Button Was Clicked!')">
    <button @click="$event.target.dispatchEvent(new CustomEvent('foo', { bubbles: true }))">...</button>
</div>
```
When the button is clicked, the `@foo` listener will be called.
Because the `.dispatchEvent` API is verbose, Alpine offers a `$dispatch` helper to simplify things.
Here's the same component re-written with the `$dispatch` magic property.
```alpine
<div x-data @foo="alert('Button Was Clicked!')">
    <button @click="$dispatch('foo')">...</button>
</div>
```
## Modifiers
Alpine offers a number of directive modifiers to customize the behavior of your event listeners.
### .prevent
`.prevent` is the equivalent of calling `.preventDefault()` inside a listener on the browser event object.
```alpine
<form @submit.prevent="console.log('submitted')" action="/foo">
    <button>Submit</button>
</form>
```
In the above example, with the `.prevent`, clicking the button will NOT submit the form to the `/foo` endpoint. Instead, Alpine's listener will handle it and "prevent" the event from being handled any further.
### .stop
Similar to `.prevent`, `.stop` is the equivalent of calling `.stopPropagation()` inside a listener on the browser event object.
```alpine
<div @click="console.log('I will not get logged')">
    <button @click.stop>Click Me</button>
</div>
```
In the above example, clicking the button WON'T log the message. This is because we are stopping the propagation of the event immediately and not allowing it to "bubble" up to the `<div>` with the `@click` listener on it.
### .outside
`.outside` is a convenience helper for listening for a click outside of the element it is attached to. Here's a simple dropdown component example to demonstrate:
```alpine
<div x-data="{ open: false }">
    <button @click="open = ! open">Toggle</button>
    <div x-show="open" @click.outside="open = false">
        Contents...
    </div>
</div>
```
In the above example, after showing the dropdown contents by clicking the "Toggle" button, you can close the dropdown by clicking anywhere on the page outside the content.
This is because `.outside` is listening for clicks that DON'T originate from the element it's registered on.
> It's worth noting that the `.outside` expression will only be evaluated when the element it's registered on is visible on the page. Otherwise, there would be nasty race conditions where clicking the "Toggle" button would also fire the `@click.outside` handler when it is not visible.
### .window
When the `.window` modifier is present, Alpine will register the event listener on the root `window` object on the page instead of the element itself.
```alpine
<div @keyup.escape.window="...">...</div>
```
The above snippet will listen for the "escape" key to be pressed ANYWHERE on the page.
Adding `.window` to listeners is extremely useful for these sorts of cases where a small part of your markup is concerned with events that take place on the entire page.
### .document
`.document` works similarly to `.window` only it registers listeners on the `document` global, instead of the `window` global.
### .once
By adding `.once` to a listener, you are ensuring that the handler is only called ONCE.
```alpine
<button @click.once="console.log('I will only log once')">...</button>
```
### .debounce
Sometimes it is useful to "debounce" an event handler so that it only is called after a certain period of inactivity (250 milliseconds by default).
For example if you have a search field that fires network requests as the user types into it, adding a debounce will prevent the network requests from firing on every single keystroke.
```alpine
<input @input.debounce="fetchResults">
```
Now, instead of calling `fetchResults` after every keystroke, `fetchResults` will only be called after 250 milliseconds of no keystrokes.
If you wish to lengthen or shorten the debounce time, you can do so by trailing a duration after the `.debounce` modifier like so:
```alpine
<input @input.debounce.500ms="fetchResults">
```
Now, `fetchResults` will only be called after 500 milliseconds of inactivity.
### .throttle
`.throttle` is similar to `.debounce` except it will release a handler call every 250 milliseconds instead of deferring it indefinitely.
This is useful for cases where there may be repeated and prolonged event firing and using `.debounce` won't work because you want to still handle the event every so often.
For example:
```alpine
<div @scroll.window.throttle="handleScroll">...</div>
```
The above example is a great use case of throttling. Without `.throttle`, the `handleScroll` method would be fired hundreds of times as the user scrolls down a page. This can really slow down a site. By adding `.throttle`, we are ensuring that `handleScroll` only gets called every 250 milliseconds.
> Fun Fact: This exact strategy is used on this very documentation site to update the currently highlighted section in the right sidebar.
Just like with `.debounce`, you can add a custom duration to your throttled event:
```alpine
<div @scroll.window.throttle.750ms="handleScroll">...</div>
```
Now, `handleScroll` will only be called every 750 milliseconds.
### .self
By adding `.self` to an event listener, you are ensuring that the event originated on the element it is declared on, and not from a child element.
```alpine
<button @click.self="handleClick">
    Click Me
    <img src="...">
</button>
```
In the above example, we have an `<img>` tag inside the `<button>` tag. Normally, any click originating within the `<button>` element (like on `<img>` for example), would be picked up by a `@click` listener on the button.
However, in this case, because we've added a `.self`, only clicking the button itself will call `handleClick`. Only clicks originating on the `<img>` element will not be handled.
### .camel
```alpine
<div @custom-event.camel="handleCustomEvent">
    ...
</div>
```
Sometimes you may want to listen for camelCased events such as `customEvent` in our example. Because camelCasing inside HTML attributes is not supported, adding the `.camel` modifier is necessary for Alpine to camelCase the event name internally.
By adding `.camel` in the above example, Alpine is now listening for `customEvent` instead of `custom-event`.
### .dot
```alpine
<div @custom-event.dot="handleCustomEvent">
    ...
</div>
```
Similar to the `.camelCase` modifier there may be situations where you want to listen for events that have dots in their name (like `custom.event`). Since dots within the event name are reserved by Alpine you need to write them with dashes and add the `.dot` modifier.
In the code example above `custom-event.dot` will correspond to the event name `custom.event`.
### .passive
Browsers optimize scrolling on pages to be fast and smooth even when JavaScript is being executed on the page. However, improperly implemented touch and wheel listeners can block this optimization and cause poor site performance.
If you are listening for touch events, it's important to add `.passive` to your listeners to not block scroll performance.
```alpine
<div @touchstart.passive="...">...</div>
```
### .capture
Add this modifier if you want to execute this listener in the event's capturing phase, e.g. before the event bubbles from the target element up the DOM.
```alpine
<div @click.capture="console.log('I will log first')">
    <button @click="console.log('I will log second')"></button>
</div>
```
##### directives-ref
# x-ref
`x-ref` in combination with `$refs` is a useful utility for easily accessing DOM elements directly. It's most useful as a replacement for APIs like `getElementById` and `querySelector`.
```alpine
<button @click="$refs.text.remove()">Remove Text</button>
<span x-ref="text">Hello 👋</span>
```
<div class="demo">
    <div x-data>
        <button @click="$refs.text.remove()">Remove Text</button>
        <div class="pt-4" x-ref="text">Hello 👋</div>
    </div>
</div>
> Despite not being included in the above snippet, `x-ref` cannot be used if no parent element has `x-data` defined. 
##### directives-show
# x-show
`x-show` is one of the most useful and powerful directives in Alpine. It provides an expressive way to show and hide DOM elements.
Here's an example of a simple dropdown component using `x-show`.
```alpine
<div x-data="{ open: false }">
    <button x-on:click="open = ! open">Toggle Dropdown</button>
    <div x-show="open">
        Dropdown Contents...
    </div>
</div>
```
When the "Toggle Dropdown" button is clicked, the dropdown will show and hide accordingly.
> If the "default" state of an `x-show` on page load is "false", you may want to use `x-cloak` on the page to avoid "page flicker" (The effect that happens when the browser renders your content before Alpine is finished initializing and hiding it.) You can learn more about `x-cloak` in its documentation.
## With transitions
If you want to apply smooth transitions to the `x-show` behavior, you can use it in conjunction with `x-transition`. You can learn more about that directive [here](/directives/transition), but here's a quick example of the same component as above, just with transitions applied.
```alpine
<div x-data="{ open: false }">
    <button x-on:click="open = ! open">Toggle Dropdown</button>
    <div x-show="open" x-transition>
        Dropdown Contents...
    </div>
</div>
```
## Using the important modifier
Sometimes you need to apply a little more force to actually hide an element. In cases where a CSS selector applies the `display` property with the `!important` flag, it will take precedence over the inline style set by Alpine.
In these cases you may use the `.important` modifier to set the inline style to `display: none !important`.
```alpine
<div x-data="{ open: false }">
    <button x-on:click="open = ! open">Toggle Dropdown</button>
    <div x-show.important="open">
        Dropdown Contents...
    </div>
</div>
```
##### directives-teleport
# x-teleport
The `x-teleport` directive allows you to transport part of your Alpine template to another part of the DOM on the page entirely.
This is useful for things like modals (especially nesting them), where it's helpful to break out of the z-index of the current Alpine component.
## x-teleport
By attaching `x-teleport` to a `<template>` element, you are telling Alpine to "append" that element to the provided selector.
> The `x-teleport` selector can be any string you would normally pass into something like `document.querySelector`. It will find the first element that matches, be it a tag name (`body`), class name (`.my-class`), ID (`#my-id`), or any other valid CSS selector.
Here's a contrived modal example:
```alpine
<body>
    <div x-data="{ open: false }">
        <button @click="open = ! open">Toggle Modal</button>
        <template x-teleport="body">
            <div x-show="open">
                Modal contents...
            </div>
        </template>
    </div>
    <div>Some other content placed AFTER the modal markup.</div>
    ...
</body>
```
<div class="demo" x-ref="root" id="modal2">
    <div x-data="{ open: false }">
        <button @click="open = ! open">Toggle Modal</button>
        <template x-teleport="#modal2">
            <div x-show="open">
                Modal contents...
            </div>
        </template>
    </div>
    <div class="py-4">Some other content placed AFTER the modal markup.</div>
</div>
Notice how when toggling the modal, the actual modal contents show up AFTER the "Some other content..." element? This is because when Alpine is initializing, it sees `x-teleport="body"` and appends and initializes that element to the provided element selector.
## Forwarding events
Alpine tries its best to make the experience of teleporting seamless. Anything you would normally do in a template, you should be able to do inside an `x-teleport` template. Teleported content can access the normal Alpine scope of the component as well as other features like `$refs`, `$root`, etc...
However, native DOM events have no concept of teleportation, so if, for example, you trigger a "click" event from inside a teleported element, that event will bubble up the DOM tree as it normally would.
To make this experience more seamless, you can "forward" events by simply registering event listeners on the `<template x-teleport...>` element itself like so:
```alpine
<div x-data="{ open: false }">
    <button @click="open = ! open">Toggle Modal</button>
    <template x-teleport="body" @click="open = false">
        <div x-show="open">
            Modal contents...
            (click to close)
        </div>
    </template>
</div>
```
<div class="demo" x-ref="root" id="modal3">
    <div x-data="{ open: false }">
        <button @click="open = ! open">Toggle Modal</button>
        <template x-teleport="#modal3" @click="open = false">
            <div x-show="open">
                Modal contents...
                <div>(click to close)</div>
            </div>
        </template>
    </div>
</div>
Notice how we are now able to listen for events dispatched from within the teleported element from outside the `<template>` element itself?
Alpine does this by looking for event listeners registered on `<template x-teleport...>` and stops those events from propagating past the live, teleported, DOM element. Then, it creates a copy of that event and re-dispatches it from `<template x-teleport...>`.
## Nesting
Teleporting is especially helpful if you are trying to nest one modal within another. Alpine makes it simple to do so:
```alpine
<div x-data="{ open: false }">
    <button @click="open = ! open">Toggle Modal</button>
    <template x-teleport="body">
        <div x-show="open">
            Modal contents...
            <div x-data="{ open: false }">
                <button @click="open = ! open">Toggle Nested Modal</button>
                <template x-teleport="body">
                    <div x-show="open">
                        Nested modal contents...
                    </div>
                </template>
            </div>
        </div>
    </template>
</div>
```
<div class="demo" x-ref="root" id="modal4">
    <div x-data="{ open: false }">
        <button @click="open = ! open">Toggle Modal</button>
        <template x-teleport="#modal4">
            <div x-show="open">
                <div class="py-4">Modal contents...</div>
                <div x-data="{ open: false }">
                    <button @click="open = ! open">Toggle Nested Modal</button>
                    <template x-teleport="#modal4">
                        <div class="pt-4" x-show="open">
                            Nested modal contents...
                        </div>
                    </template>
                </div>
            </div>
        </template>
    </div>
    <template x-teleport-target="modals3"></template>
</div>
After toggling "on" both modals, they are authored as children, but will be rendered as sibling elements on the page, not within one another.
##### directives-text
# x-text
`x-text` sets the text content of an element to the result of a given expression.
Here's a basic example of using `x-text` to display a user's username.
```alpine
<div x-data="{ username: 'calebporzio' }">
    Username: <strong x-text="username"></strong>
</div>
```
<div class="demo">
    <div x-data="{ username: 'calebporzio' }">
        Username: <strong x-text="username"></strong>
    </div>
</div>
Now the `<strong>` tag's inner text content will be set to "calebporzio".
##### directives-transition
# x-transition
Alpine provides a robust transitions utility out of the box. With a few `x-transition` directives, you can create smooth transitions between when an element is shown or hidden.
There are two primary ways to handle transitions in Alpine:
* [The Transition Helper](#the-transition-helper)
* [Applying CSS Classes](#applying-css-classes)
## The transition helper
The simplest way to achieve a transition using Alpine is by adding `x-transition` to an element with `x-show` on it. For example:
```alpine
<div x-data="{ open: false }">
    <button @click="open = ! open">Toggle</button>
    <div x-show="open" x-transition>
        Hello 👋
    </div>
</div>
```
<div class="demo">
    <div x-data="{ open: false }">
        <button @click="open = ! open">Toggle</button>
        <div x-show="open" x-transition>
            Hello 👋
        </div>
    </div>
</div>
As you can see, by default, `x-transition` applies pleasant transition defaults to fade and scale the revealing element.
You can override these defaults with modifiers attached to `x-transition`. Let's take a look at those.
### Customizing duration
Initially, the duration is set to be 150 milliseconds when entering, and 75 milliseconds when leaving.
You can configure the duration you want for a transition with the `.duration` modifier:
```alpine
<div ... x-transition.duration.500ms>
```
The above `<div>` will transition for 500 milliseconds when entering, and 500 milliseconds when leaving.
If you wish to customize the durations specifically for entering and leaving, you can do that like so:
```alpine
<div ...
    x-transition:enter.duration.500ms
    x-transition:leave.duration.400ms
>
```
> Despite not being included in the above snippet, `x-transition` cannot be used if no parent element has `x-data` defined. 
### Customizing delay
You can delay a transition using the `.delay` modifier like so:
```alpine
<div ... x-transition.delay.50ms>
```
The above example will delay the transition and in and out of the element by 50 milliseconds.
### Customizing opacity
By default, Alpine's `x-transition` applies both a scale and opacity transition to achieve a "fade" effect.
If you wish to only apply the opacity transition (no scale), you can accomplish that like so:
```alpine
<div ... x-transition.opacity>
```
### Customizing scale
Similar to the `.opacity` modifier, you can configure `x-transition` to ONLY scale (and not transition opacity as well) like so:
```alpine
<div ... x-transition.scale>
```
The `.scale` modifier also offers the ability to configure its scale values AND its origin values:
```alpine
<div ... x-transition.scale.80>
```
The above snippet will scale the element up and down by 80%.
Again, you may customize these values separately for enter and leaving transitions like so:
```alpine
<div ...
    x-transition:enter.scale.80
    x-transition:leave.scale.90
>
```
To customize the origin of the scale transition, you can use the `.origin` modifier:
```alpine
<div ... x-transition.scale.origin.top>
```
Now the scale will be applied using the top of the element as the origin, instead of the center by default.
Like you may have guessed, the possible values for this customization are: `top`, `bottom`, `left`, and `right`.
If you wish, you can also combine two origin values. For example, if you want the origin of the scale to be "top right", you can use: `.origin.top.right` as the modifier.
## Applying CSS classes
For direct control over exactly what goes into your transitions, you can apply CSS classes at different stages of the transition.
> The following examples use [TailwindCSS](https://tailwindcss.com/docs/transition-property) utility classes.
```alpine
<div x-data="{ open: false }">
    <button @click="open = ! open">Toggle</button>
    <div
        x-show="open"
        x-transition:enter="transition ease-out duration-300"
        x-transition:enter-start="opacity-0 scale-90"
        x-transition:enter-end="opacity-100 scale-100"
        x-transition:leave="transition ease-in duration-300"
        x-transition:leave-start="opacity-100 scale-100"
        x-transition:leave-end="opacity-0 scale-90"
    >Hello 👋</div>
</div>
```
<div class="demo">
    <div x-data="{ open: false }">
    <button @click="open = ! open">Toggle</button>
    <div
        x-show="open"
        x-transition:enter="transition ease-out duration-300"
        x-transition:enter-start="opacity-0 transform scale-90"
        x-transition:enter-end="opacity-100 transform scale-100"
        x-transition:leave="transition ease-in duration-300"
        x-transition:leave-start="opacity-100 transform scale-100"
        x-transition:leave-end="opacity-0 transform scale-90"
    >Hello 👋</div>
</div>
</div>
| Directive      | Description |
| ---            | --- |
| `:enter`       | Applied during the entire entering phase. |
| `:enter-start` | Added before element is inserted, removed one frame after element is inserted. |
| `:enter-end`   | Added one frame after element is inserted (at the same time `enter-start` is removed), removed when transition/animation finishes.
| `:leave`       | Applied during the entire leaving phase. |
| `:leave-start` | Added immediately when a leaving transition is triggered, removed after one frame. |
| `:leave-end`   | Added one frame after a leaving transition is triggered (at the same time `leave-start` is removed), removed when the transition/animation finishes.
##### essentials-events
# Events
Alpine makes it simple to listen for browser events and react to them.
## Listening for simple events
By using `x-on`, you can listen for browser events that are dispatched on or within an element.
Here's a basic example of listening for a click on a button:
```alpine
<button x-on:click="console.log('clicked')">...</button>
```
As an alternative, you can use the event shorthand syntax if you prefer: `@`. Here's the same example as before, but using the shorthand syntax (which we'll be using from now on):
```alpine
<button @click="...">...</button>
```
In addition to `click`, you can listen for any browser event by name. For example: `@mouseenter`, `@keyup`, etc... are all valid syntax.
## Listening for specific keys
Let's say you wanted to listen for the `enter` key to be pressed inside an `<input>` element. Alpine makes this easy by adding the `.enter` like so:
```alpine
<input @keyup.enter="...">
```
You can even combine key modifiers to listen for key combinations like pressing `enter` while holding `shift`:
```alpine
<input @keyup.shift.enter="...">
```
## Preventing default
When reacting to browser events, it is often necessary to "prevent default" (prevent the default behavior of the browser event).
For example, if you want to listen for a form submission but prevent the browser from submitting a form request, you can use `.prevent`:
```alpine
<form @submit.prevent="...">...</form>
```
You can also apply `.stop` to achieve the equivalent of `event.stopPropagation()`.
## Accessing the event object
Sometimes you may want to access the native browser event object inside your own code. To make this easy, Alpine automatically injects an `$event` magic variable:
```alpine
<button @click="$event.target.remove()">Remove Me</button>
```
## Dispatching custom events
In addition to listening for browser events, you can dispatch them as well. This is extremely useful for communicating with other Alpine components or triggering events in tools outside of Alpine itself.
Alpine exposes a magic helper called `$dispatch` for this:
```alpine
<div @foo="console.log('foo was dispatched')">
    <button @click="$dispatch('foo')"></button>
</div>
```
As you can see, when the button is clicked, Alpine will dispatch a browser event called "foo", and our `@foo` listener on the `<div>` will pick it up and react to it.
## Listening for events on window
Because of the nature of events in the browser, it is sometimes useful to listen to events on the top-level window object.
This allows you to communicate across components completely like the following example:
```alpine
<div x-data>
    <button @click="$dispatch('foo')"></button>
</div>
<div x-data @foo.window="console.log('foo was dispatched')">...</div>
```
In the above example, if we click the button in the first component, Alpine will dispatch the "foo" event. Because of the way events work in the browser, they "bubble" up through parent elements all the way to the top-level "window".
Now, because in our second component we are listening for "foo" on the window (with `.window`), when the button is clicked, this listener will pick it up and log the "foo was dispatched" message.
##### essentials-lifecycle
# Lifecycle
Alpine has a handful of different techniques for hooking into different parts of its lifecycle. Let's go through the most useful ones to familiarize yourself with:
## Element initialization
Another extremely useful lifecycle hook in Alpine is the `x-init` directive.
`x-init` can be added to any element on a page and will execute any JavaScript you call inside it when Alpine begins initializing that element.
```alpine
<button x-init="console.log('Im initing')">
```
In addition to the directive, Alpine will automatically call any `init()` methods stored on a data object. For example:
```js
Alpine.data('dropdown', () => ({
    init() {
        // I get called before the element using this data initializes.
    }
}))
```
## After a state change
Alpine allows you to execute code when a piece of data (state) changes. It offers two different APIs for such a task: `$watch` and `x-effect`.
### `$watch`
```alpine
<div x-data="{ open: false }" x-init="$watch('open', value => console.log(value))">
```
As you can see above, `$watch` allows you to hook into data changes using a dot-notation key. When that piece of data changes, Alpine will call the passed callback and pass it the new value. along with the old value before the change.
### `x-effect`
`x-effect` uses the same mechanism under the hood as `$watch` but has very different usage.
Instead of specifying which data key you wish to watch, `x-effect` will call the provided code and intelligently look for any Alpine data used within it. Now when one of those pieces of data changes, the `x-effect` expression will be re-run.
Here's the same bit of code from the `$watch` example rewritten using `x-effect`:
```alpine
<div x-data="{ open: false }" x-effect="console.log(open)">
```
Now, this expression will be called right away, and re-called every time `open` is updated.
The two main behavioral differences with this approach are:
1. The provided code will be run right away AND when data changes (`$watch` is "lazy" -- won't run until the first data change)
2. No knowledge of the previous value. (The callback provided to `$watch` receives both the new value AND the old one)
## Alpine initialization
### `alpine:init`
Ensuring a bit of code executes after Alpine is loaded, but BEFORE it initializes itself on the page is a necessary task.
This hook allows you to register custom data, directives, magics, etc. before Alpine does its thing on a page.
You can hook into this point in the lifecycle by listening for an event that Alpine dispatches called: `alpine:init`
```js
document.addEventListener('alpine:init', () => {
    Alpine.data(...)
})
```
### `alpine:initialized`
Alpine also offers a hook that you can use to execute code AFTER it's done initializing called `alpine:initialized`:
```js
document.addEventListener('alpine:initialized', () => {
    //
})
```
##### essentials-state
# State
State (JavaScript data that Alpine watches for changes) is at the core of everything you do in Alpine. You can provide local data to a chunk of HTML, or make it globally available for use anywhere on a page using `x-data` or `Alpine.store()` respectively.
## Local state
Alpine allows you to declare an HTML block's state in a single `x-data` attribute without ever leaving your markup.
Here's a basic example:
```alpine
<div x-data="{ open: false }">
    ...
</div>
```
Now any other Alpine syntax on or within this element will be able to access `open`. And like you'd guess, when `open` changes for any reason, everything that depends on it will react automatically.
### Nesting data
Data is nestable in Alpine. For example, if you have two elements with Alpine data attached (one inside the other), you can access the parent's data from inside the child element.
```alpine
<div x-data="{ open: false }">
    <div x-data="{ label: 'Content:' }">
        <span x-text="label"></span>
        <span x-show="open"></span>
    </div>
</div>
```
This is similar to scoping in JavaScript itself (code within a function can access variables declared outside that function.)
Like you may have guessed, if the child has a data property matching the name of a parent's property, the child property will take precedence.
### Single-element data
Although this may seem obvious to some, it's worth mentioning that Alpine data can be used within the same element. For example:
```alpine
<button x-data="{ label: 'Click Here' }" x-text="label"></button>
```
### Data-less Alpine
Sometimes you may want to use Alpine functionality, but don't need any reactive data. In these cases, you can opt out of passing an expression to `x-data` entirely. For example:
```alpine
<button x-data @click="alert('I\'ve been clicked!')">Click Me</button>
```
### Re-usable data
When using Alpine, you may find the need to re-use a chunk of data and/or its corresponding template.
If you are using a backend framework like Rails or Laravel, Alpine first recommends that you extract the entire block of HTML into a template partial or include.
If for some reason that isn't ideal for you or you're not in a back-end templating environment, Alpine allows you to globally register and re-use the data portion of a component using `Alpine.data(...)`.
```js
Alpine.data('dropdown', () => ({
    open: false,
    toggle() {
        this.open = ! this.open
    }
}))
```
Now that you've registered the "dropdown" data, you can use it inside your markup in as many places as you like:
```alpine
<div x-data="dropdown">
    <button @click="toggle">Expand</button>
    <span x-show="open">Content...</span>
</div>
<div x-data="dropdown">
    <button @click="toggle">Expand</button>
    <span x-show="open">Some Other Content...</span>
</div>
```
## Global state
If you wish to make some data available to every component on the page, you can do so using Alpine's "global store" feature.
You can register a store using `Alpine.store(...)`, and reference one with the magic `$store()` method.
Let's look at a simple example. First we'll register the store globally:
```js
Alpine.store('tabs', {
    current: 'first',
    items: ['first', 'second', 'third'],
})
```
Now we can access or modify its data from anywhere on our page:
```alpine
<div x-data>
    <template x-for="tab in $store.tabs.items">
        ...
    </template>
</div>
<div x-data>
    <button @click="$store.tabs.current = 'first'">First Tab</button>
    <button @click="$store.tabs.current = 'second'">Second Tab</button>
    <button @click="$store.tabs.current = 'third'">Third Tab</button>
</div>
```
##### essentials-templating
# Templating
Alpine offers a handful of useful directives for manipulating the DOM on a web page.
Let's cover a few of the basic templating directives here, but be sure to look through the available directives in the sidebar for an exhaustive list.
## Text content
Alpine makes it easy to control the text content of an element with the `x-text` directive.
```alpine
<div x-data="{ title: 'Start Here' }">
    <h1 x-text="title"></h1>
</div>
```
<div x-data="{ title: 'Start Here' }" class="demo">
    <strong x-text="title"></strong>
</div>
Now, Alpine will set the text content of the `<h1>` with the value of `title` ("Start Here"). When `title` changes, so will the contents of `<h1>`.
Like all directives in Alpine, you can use any JavaScript expression you like. For example:
```alpine
<span x-text="1 + 2"></span>
```
<div class="demo" x-data>
    <span x-text="1 + 2"></span>
</div>
The `<span>` will now contain the sum of "1" and "2".
## Toggling elements
Toggling elements is a common need in web pages and applications. Dropdowns, modals, dialogues, "show-more"s, etc... are all good examples.
Alpine offers the `x-show` and `x-if` directives for toggling elements on a page.
### `x-show`
Here's a simple toggle component using `x-show`.
```alpine
<div x-data="{ open: false }">
    <button @click="open = ! open">Expand</button>
    <div x-show="open">
        Content...
    </div>
</div>
```
<div x-data="{ open: false }" class="demo">
    <button @click="open = ! open" :aria-pressed="open">Expand</button>
    <div x-show="open">
        Content...
    </div>
</div>
Now the entire `<div>` containing the contents will be shown and hidden based on the value of `open`.
Under the hood, Alpine adds the CSS property `display: none;` to the element when it should be hidden.
This works well for most cases, but sometimes you may want to completely add and remove the element from the DOM entirely. This is what `x-if` is for.
### `x-if`
Here is the same toggle from before, but this time using `x-if` instead of `x-show`.
```alpine
<div x-data="{ open: false }">
    <button @click="open = ! open">Expand</button>
    <template x-if="open">
        <div>
            Content...
        </div>
    </template>
</div>
```
<div x-data="{ open: false }" class="demo">
    <button @click="open = ! open" :aria-pressed="open">Expand</button>
    <template x-if="open">
        <div>
            Content...
        </div>
    </template>
</div>
Notice that `x-if` must be declared on a `<template>` tag. This is so that Alpine can leverage the existing browser behavior of the `<template>` element and use it as the source of the target `<div>` to be added and removed from the page.
When `open` is true, Alpine will append the `<div>` to the `<template>` tag, and remove it when `open` is false.
## Toggling with transitions
Alpine makes it simple to smoothly transition between "shown" and "hidden" states using the `x-transition` directive.
> `x-transition` only works with `x-show`, not with `x-if`.
Here is, again, the simple toggle example, but this time with transitions applied:
```alpine
<div x-data="{ open: false }">
    <button @click="open = ! open">Expands</button>
    <div x-show="open" x-transition>
        Content...
    </div>
</div>
```
<div x-data="{ open: false }" class="demo">
    <button @click="open = ! open">Expands</button>
    <div class="flex">
        <div x-show="open" x-transition style="will-change: transform;">
            Content...
        </div>
    </div>
</div>
Let's zoom in on the portion of the template dealing with transitions:
```alpine
<div x-show="open" x-transition>
```
`x-transition` by itself will apply sensible default transitions (fade and scale) to the toggle.
There are two ways to customize these transitions:
* Transition helpers
* Transition CSS classes.
Let's take a look at each of these approaches:
### Transition helpers
Let's say you wanted to make the duration of the transition longer, you can manually specify that using the `.duration` modifier like so:
```alpine
<div x-show="open" x-transition.duration.500ms>
```
<div x-data="{ open: false }" class="demo">
    <button @click="open = ! open">Expands</button>
    <div class="flex">
        <div x-show="open" x-transition.duration.500ms style="will-change: transform;">
            Content...
        </div>
    </div>
</div>
Now the transition will last 500 milliseconds.
If you want to specify different values for in and out transitions, you can use `x-transition:enter` and `x-transition:leave`:
```alpine
<div
    x-show="open"
    x-transition:enter.duration.500ms
    x-transition:leave.duration.1000ms
>
```
<div x-data="{ open: false }" class="demo">
    <button @click="open = ! open">Expands</button>
    <div class="flex">
        <div x-show="open" x-transition:enter.duration.500ms x-transition:leave.duration.1000ms style="will-change: transform;">
            Content...
        </div>
    </div>
</div>
Additionally, you can add either `.opacity` or `.scale` to only transition that property. For example:
```alpine
<div x-show="open" x-transition.opacity>
```
<div x-data="{ open: false }" class="demo">
    <button @click="open = ! open">Expands</button>
    <div class="flex">
        <div x-show="open" x-transition:enter.opacity.duration.500 x-transition:leave.opacity.duration.250>
            Content...
        </div>
    </div>
</div>
### Transition classes
If you need more fine-grained control over the transitions in your application, you can apply specific CSS classes at specific phases of the transition using the following syntax (this example uses [Tailwind CSS](https://tailwindcss.com/)):
```alpine
<div
    x-show="open"
    x-transition:enter="transition ease-out duration-300"
    x-transition:enter-start="opacity-0 transform scale-90"
    x-transition:enter-end="opacity-100 transform scale-100"
    x-transition:leave="transition ease-in duration-300"
    x-transition:leave-start="opacity-100 transform scale-100"
    x-transition:leave-end="opacity-0 transform scale-90"
>...</div>
```
<div x-data="{ open: false }" class="demo">
    <button @click="open = ! open">Expands</button>
    <div class="flex">
        <div
            x-show="open"
            x-transition:enter="transition ease-out duration-300"
            x-transition:enter-start="opacity-0 transform scale-90"
            x-transition:enter-end="opacity-100 transform scale-100"
            x-transition:leave="transition ease-in duration-300"
            x-transition:leave-start="opacity-100 transform scale-100"
            x-transition:leave-end="opacity-0 transform scale-90"
            style="will-change: transform"
        >
            Content...
        </div>
    </div>
</div>
## Binding attributes
You can add HTML attributes like `class`, `style`, `disabled`, etc... to elements in Alpine using the `x-bind` directive.
Here is an example of a dynamically bound `class` attribute:
```alpine
<button
    x-data="{ red: false }"
    x-bind:class="red ? 'bg-red' : ''"
    @click="red = ! red"
>
    Toggle Red
</button>
```
<div class="demo">
    <button
        x-data="{ red: false }"
        x-bind:style="red && 'background: red'"
        @click="red = ! red"
    >
        Toggle Red
    </button>
</div>
As a shortcut, you can leave out the `x-bind` and use the shorthand `:` syntax directly:
```alpine
<button ... :class="red ? 'bg-red' : ''">
```
Toggling classes on and off based on data inside Alpine is a common need. Here's an example of toggling a class using Alpine's `class` binding object syntax: (Note: this syntax is only available for `class` attributes)
```alpine
<div x-data="{ open: true }">
    <span :class="{ 'hidden': ! open }">...</span>
</div>
```
Now the `hidden` class will be added to the element if `open` is false, and removed if `open` is true.
## Looping elements
Alpine allows for iterating parts of your template based on JavaScript data using the `x-for` directive. Here is a simple example:
```alpine
<div x-data="{ statuses: ['open', 'closed', 'archived'] }">
    <template x-for="status in statuses">
        <div x-text="status"></div>
    </template>
</div>
```
<div x-data="{ statuses: ['open', 'closed', 'archived'] }" class="demo">
    <template x-for="status in statuses">
        <div x-text="status"></div>
    </template>
</div>
Similar to `x-if`, `x-for` must be applied to a `<template>` tag. Internally, Alpine will append the contents of `<template>` tag for every iteration in the loop.
As you can see the new `status` variable is available in the scope of the iterated templates.
## Inner HTML
Alpine makes it easy to control the HTML content of an element with the `x-html` directive.
```alpine
<div x-data="{ title: '<h1>Start Here</h1>' }">
    <div x-html="title"></div>
</div>
```
<div x-data="{ title: '<h1>Start Here</h1>' }" class="demo">
    <div x-html="title"></div>
</div>
Now, Alpine will set the text content of the `<div>` with the element `<h1>Start Here</h1>`. When `title` changes, so will the contents of `<h1>`.
> ⚠️ Only use on trusted content and never on user-provided content. ⚠️
> Dynamically rendering HTML from third parties can easily lead to XSS vulnerabilities.
##### globals-alpine-bind
# Alpine.bind
`Alpine.bind(...)` provides a way to re-use [`x-bind`](/directives/bind#bind-directives) objects within your application.
Here's a simple example. Rather than binding attributes manually with Alpine:
```alpine
<button type="button" @click="doSomething()" :disabled="shouldDisable"></button>
```
You can bundle these attributes up into a reusable object and use `x-bind` to bind to that:
```alpine
<button x-bind="SomeButton"></button>
<script>
    document.addEventListener('alpine:init', () => {
        Alpine.bind('SomeButton', () => ({
            type: 'button',
            '@click'() {
                this.doSomething()
            },
            ':disabled'() {
                return this.shouldDisable
            },
        }))
    })
</script>
```
##### globals-alpine-data
# Alpine.data
`Alpine.data(...)` provides a way to re-use `x-data` contexts within your application.
Here's a contrived `dropdown` component for example:
```alpine
<div x-data="dropdown">
    <button @click="toggle">...</button>
    <div x-show="open">...</div>
</div>
<script>
    document.addEventListener('alpine:init', () => {
        Alpine.data('dropdown', () => ({
            open: false,
            toggle() {
                this.open = ! this.open
            }
        }))
    })
</script>
```
As you can see we've extracted the properties and methods we would usually define directly inside `x-data` into a separate Alpine component object.
## Registering from a bundle
If you've chosen to use a build step for your Alpine code, you should register your components in the following way:
```js
import Alpine from 'alpinejs'
import dropdown from './dropdown.js'
Alpine.data('dropdown', dropdown)
Alpine.start()
```
This assumes you have a file called `dropdown.js` with the following contents:
```js
export default () => ({
    open: false,
    toggle() {
        this.open = ! this.open
    }
})
```
## Initial parameters
In addition to referencing `Alpine.data` providers by their name plainly (like `x-data="dropdown"`), you can also reference them as functions (`x-data="dropdown()"`). By calling them as functions directly, you can pass in additional parameters to be used when creating the initial data object like so:
```alpine
<div x-data="dropdown(true)">
```
```js
Alpine.data('dropdown', (initialOpenState = false) => ({
    open: initialOpenState
}))
```
Now, you can re-use the `dropdown` object, but provide it with different parameters as you need to.
## Init functions
If your component contains an `init()` method, Alpine will automatically execute it before it renders the component. For example:
```js
Alpine.data('dropdown', () => ({
    init() {
        // This code will be executed before Alpine
        // initializes the rest of the component.
    }
}))
```
## Destroy functions
If your component contains a `destroy()` method, Alpine will automatically execute it before cleaning up the component.
A primary example for this is when registering an event handler with another library or a browser API that isn't available through Alpine.
See the following example code on how to use the `destroy()` method to clean up such a handler.
```js
Alpine.data('timer', () => ({
    timer: null,
    counter: 0,
    init() {
      // Register an event handler that references the component instance
      this.timer = setInterval(() => {
        console.log('Increased counter to', ++this.counter);
      }, 1000);
    },
    destroy() {
        // Detach the handler, avoiding memory and side-effect leakage
        clearInterval(this.timer);
    },
}))
```
An example where a component is destroyed is when using one inside an `x-if`:
```html
<span x-data="{ enabled: false }">
    <button @click.prevent="enabled = !enabled">Toggle</button>
    <template x-if="enabled">
        <span x-data="timer" x-text="counter"></span>
    </template>
</span>
```
## Using magic properties
If you want to access magic methods or properties from a component object, you can do so using the `this` context:
```js
Alpine.data('dropdown', () => ({
    open: false,
    init() {
        this.$watch('open', () => {...})
    }
}))
```
## Encapsulating directives with `x-bind`
If you wish to re-use more than just the data object of a component, you can encapsulate entire Alpine template directives using `x-bind`.
The following is an example of extracting the templating details of our previous dropdown component using `x-bind`:
```alpine
<div x-data="dropdown">
    <button x-bind="trigger"></button>
    <div x-bind="dialogue"></div>
</div>
```
```js
Alpine.data('dropdown', () => ({
    open: false,
    trigger: {
        ['@click']() {
            this.open = ! this.open
        },
    },
    dialogue: {
        ['x-show']() {
            return this.open
        },
    },
}))
```
##### globals-alpine-store
# Alpine.store
Alpine offers global state management through the `Alpine.store()` API.
## Registering A Store
You can either define an Alpine store inside of an `alpine:init` listener (in the case of including Alpine via a `<script>` tag), OR you can define it before manually calling `Alpine.start()` (in the case of importing Alpine into a build):
**From a script tag:**
```alpine
<script>
    document.addEventListener('alpine:init', () => {
        Alpine.store('darkMode', {
            on: false,
            toggle() {
                this.on = ! this.on
            }
        })
    })
</script>
```
**From a bundle:**
```js
import Alpine from 'alpinejs'
Alpine.store('darkMode', {
    on: false,
    toggle() {
        this.on = ! this.on
    }
})
Alpine.start()
```
## Accessing stores
You can access data from any store within Alpine expressions using the `$store` magic property:
```alpine
<div x-data :class="$store.darkMode.on && 'bg-black'">...</div>
```
You can also modify properties within the store and everything that depends on those properties will automatically react. For example:
```alpine
<button x-data @click="$store.darkMode.toggle()">Toggle Dark Mode</button>
```
Additionally, you can access a store externally using `Alpine.store()` by omitting the second parameter like so:
```alpine
<script>
    Alpine.store('darkMode').toggle()
</script>
```
## Initializing stores
If you provide `init()` method in an Alpine store, it will be executed right after the store is registered. This is useful for initializing any state inside the store with sensible starting values.
```alpine
<script>
    document.addEventListener('alpine:init', () => {
        Alpine.store('darkMode', {
            init() {
                this.on = window.matchMedia('(prefers-color-scheme: dark)').matches
            },
            on: false,
            toggle() {
                this.on = ! this.on
            }
        })
    })
</script>
```
Notice the newly added `init()` method in the example above. With this addition, the `on` store variable will be set to the browser's color scheme preference before Alpine renders anything on the page.
## Single-value stores
If you don't need an entire object for a store, you can set and use any kind of data as a store.
Here's the example from above but using it more simply as a boolean value:
```alpine
<button x-data @click="$store.darkMode = ! $store.darkMode">Toggle Dark Mode</button>
...
<div x-data :class="$store.darkMode && 'bg-black'">
    ...
</div>
<script>
    document.addEventListener('alpine:init', () => {
        Alpine.store('darkMode', false)
    })
</script>
```
##### magics-data
# $data
`$data` is a magic property that gives you access to the current Alpine data scope (generally provided by `x-data`).
Most of the time, you can just access Alpine data within expressions directly. for example `x-data="{ message: 'Hello Caleb!' }"` will allow you to do things like `x-text="message"`.
However, sometimes it is helpful to have an actual object that encapsulates all scope that you can pass around to other functions:
```alpine
<div x-data="{ greeting: 'Hello' }">
    <div x-data="{ name: 'Caleb' }">
        <button @click="sayHello($data)">Say Hello</button>
    </div>
</div>
<script>
    function sayHello({ greeting, name }) {
        alert(greeting + ' ' + name + '!')
    }
</script>
```
<div x-data="{ greeting: 'Hello' }" class="demo">
    <div x-data="{ name: 'Caleb' }">
        <button @click="sayHello($data)">Say Hello</button>
    </div>
</div>
<script>
    function sayHello({ greeting, name }) {
        alert(greeting + ' ' + name + '!')
    }
</script>
Now when the button is pressed, the browser will alert `Hello Caleb!` because it was passed a data object that contained all the Alpine scope of the expression that called it (`@click="..."`).
Most applications won't need this magic property, but it can be very helpful for deeper, more complicated Alpine utilities.
##### magics-dispatch
# $dispatch
`$dispatch` is a helpful shortcut for dispatching browser events.
```alpine
<div @notify="alert('Hello World!')">
    <button @click="$dispatch('notify')">
        Notify
    </button>
</div>
```
<div class="demo">
    <div x-data @notify="alert('Hello World!')">
        <button @click="$dispatch('notify')">
            Notify
        </button>
    </div>
</div>
You can also pass data along with the dispatched event if you wish. This data will be accessible as the `.detail` property of the event:
```alpine
<div @notify="alert($event.detail.message)">
    <button @click="$dispatch('notify', { message: 'Hello World!' })">
        Notify
    </button>
</div>
```
<div class="demo">
    <div x-data @notify="alert($event.detail.message)">
        <button @click="$dispatch('notify', { message: 'Hello World!' })">Notify</button>
    </div>
</div>
Under the hood, `$dispatch` is a wrapper for the more verbose API: `element.dispatchEvent(new CustomEvent(...))`
**Note on event propagation**
Notice that, because of [event bubbling](https://en.wikipedia.org/wiki/Event_bubbling), when you need to capture events dispatched from nodes that are under the same nesting hierarchy, you'll need to use the [`.window`](https://github.com/alpinejs/alpine#x-on) modifier:
**Example:**
```alpine
<div x-data>
    <span @notify="..."></span>
    <button @click="$dispatch('notify')">Notify</button>
</div>
<div x-data>
    <span @notify.window="..."></span>
    <button @click="$dispatch('notify')">Notify</button>
</div>
```
> The first example won't work because when `notify` is dispatched, it'll propagate to its common ancestor, the `div`, not its sibling, the `<span>`. The second example will work because the sibling is listening for `notify` at the `window` level, which the custom event will eventually bubble up to.
## Dispatching to other components
You can also take advantage of the previous technique to make your components talk to each other:
**Example:**
```alpine
<div
    x-data="{ title: 'Hello' }"
    @set-title.window="title = $event.detail"
>
    <h1 x-text="title"></h1>
</div>
<div x-data>
    <button @click="$dispatch('set-title', 'Hello World!')">Click me</button>
</div>
```
## Dispatching to x-model
You can also use `$dispatch()` to trigger data updates for `x-model` data bindings. For example:
```alpine
<div x-data="{ title: 'Hello' }">
    <span x-model="title">
        <button @click="$dispatch('input', 'Hello World!')">Click me</button>
        
    </span>
</div>
```
This opens up the door for making custom input components whose value can be set via `x-model`.
##### magics-el
# $el
`$el` is a magic property that can be used to retrieve the current DOM node.
```alpine
<button @click="$el.innerHTML = 'Hello World!'">Replace me with "Hello World!"</button>
```
<div class="demo">
    <div x-data>
        <button @click="$el.textContent = 'Hello World!'">Replace me with "Hello World!"</button>
    </div>
</div>
##### magics-id
# $id
`$id` is a magic property that can be used to generate an element's ID and ensure that it won't conflict with other IDs of the same name on the same page.
This utility is extremely helpful when building re-usable components (presumably in a back-end template) that might occur multiple times on a page, and make use of ID attributes.
Things like input components, modals, listboxes, etc. will all benefit from this utility.
## Basic usage
Suppose you have two input elements on a page, and you want them to have a unique ID from each other, you can do the following:
```alpine
<input type="text" :id="$id('text-input')">
<input type="text" :id="$id('text-input')">
```
As you can see, `$id` takes in a string and spits out an appended suffix that is unique on the page.
## Grouping with x-id
Now let's say you want to have those same two input elements, but this time you want `<label>` elements for each of them.
This presents a problem, you now need to be able to reference the same ID twice. One for the `<label>`'s `for` attribute, and the other for the `id` on the input.
Here is a way that you might think to accomplish this and is totally valid:
```alpine
<div x-data="{ id: $id('text-input') }">
    <label :for="id"> 
    <input type="text" :id="id"> 
</div>
<div x-data="{ id: $id('text-input') }">
    <label :for="id"> 
    <input type="text" :id="id"> 
</div>
```
This approach is fine, however, having to name and store the ID in your component scope feels cumbersome.
To accomplish this same task in a more flexible way, you can use Alpine's `x-id` directive to declare an "id scope" for a set of IDs:
```alpine
<div x-id="['text-input']">
    <label :for="$id('text-input')"> 
    <input type="text" :id="$id('text-input')"> 
</div>
<div x-id="['text-input']">
    <label :for="$id('text-input')"> 
    <input type="text" :id="$id('text-input')"> 
</div>
```
As you can see, `x-id` accepts an array of ID names. Now any usages of `$id()` within that scope, will all use the same ID. Think of them as "id groups".
## Nesting
As you might have intuited, you can freely nest these `x-id` groups, like so:
```alpine
<div x-id="['text-input']">
    <label :for="$id('text-input')"> 
    <input type="text" :id="$id('text-input')"> 
    <div x-id="['text-input']">
        <label :for="$id('text-input')"> 
        <input type="text" :id="$id('text-input')"> 
    </div>
</div>
```
## Keyed IDs (For Looping)
Sometimes, it is helpful to specify an additional suffix on the end of an ID for the purpose of identifying it within a loop.
For this, `$id()` accepts an optional second parameter that will be added as a suffix on the end of the generated ID.
A common example of this need is something like a listbox component that uses the `aria-activedescendant` attribute to tell assistive technologies which element is "active" in the list:
```alpine
<ul
    x-id="['list-item']"
    :aria-activedescendant="$id('list-item', activeItem.id)"
>
    <template x-for="item in items" :key="item.id">
        <li :id="$id('list-item', item.id)">...</li>
    </template>
</ul>
```
This is an incomplete example of a listbox, but it should still be helpful to demonstrate a scenario where you might need each ID in a group to still be unique to the page, but also be keyed within a loop so that you can reference individual IDs within that group.
##### magics-nextTick
# $nextTick
`$nextTick` is a magic property that allows you to only execute a given expression AFTER Alpine has made its reactive DOM updates. This is useful for times you want to interact with the DOM state AFTER it's reflected any data updates you've made.
```alpine
<div x-data="{ title: 'Hello' }">
    <button
        @click="
            title = 'Hello World!';
            $nextTick(() => { console.log($el.innerText) });
        "
        x-text="title"
    ></button>
</div>
```
In the above example, rather than logging "Hello" to the console, "Hello World!" will be logged because `$nextTick` was used to wait until Alpine was finished updating the DOM.
## Promises
`$nextTick` returns a promise, allowing the use of `$nextTick` to pause an async function until after pending dom updates. When used like this, `$nextTick` also does not require an argument to be passed.
```alpine
<div x-data="{ title: 'Hello' }">
    <button
        @click="
            title = 'Hello World!';
            await $nextTick();
            console.log($el.innerText);
        "
        x-text="title"
    ></button>
</div>
```
##### magics-refs
# $refs
`$refs` is a magic property that can be used to retrieve DOM elements marked with `x-ref` inside the component. This is useful when you need to manually manipulate DOM elements. It's often used as a more succinct, scoped, alternative to `document.querySelector`.
```alpine
<button @click="$refs.text.remove()">Remove Text</button>
<span x-ref="text">Hello 👋</span>
```
<div class="demo">
    <div x-data>
        <button @click="$refs.text.remove()">Remove Text</button>
        <div class="pt-4" x-ref="text">Hello 👋</div>
    </div>
</div>
Now, when the `<button>` is pressed, the `<span>` will be removed.
### Limitations
In V2 it was possible to bind `$refs` to elements dynamically, like seen below:
```alpine
<template x-for="item in items" :key="item.id" >
    <div :x-ref="item.name">
    some content ...
    </div>
</template>
```
However, in V3, `$refs` can only be accessed for elements that are created statically. So for the example above: if you were expecting the value of `item.name` inside of `$refs` to be something like *Batteries*, you should be aware that `$refs` will actually contain the literal string `'item.name'` and not *Batteries*.
##### magics-root
# $root
`$root` is a magic property that can be used to retrieve the root element of any Alpine component. In other words the closest element up the DOM tree that contains `x-data`.
```alpine
<div x-data data-message="Hello World!">
    <button @click="alert($root.dataset.message)">Say Hi</button>
</div>
```
<div x-data data-message="Hello World!" class="demo">
    <button @click="alert($root.dataset.message)">Say Hi</button>
</div>
##### magics-store
# $store
You can use `$store` to conveniently access global Alpine stores registered using [`Alpine.store(...)`](/globals/alpine-store). For example:
```alpine
<button x-data @click="$store.darkMode.toggle()">Toggle Dark Mode</button>
...
<div x-data :class="$store.darkMode.on && 'bg-black'">
    ...
</div>
<script>
    document.addEventListener('alpine:init', () => {
        Alpine.store('darkMode', {
            on: false,
            toggle() {
                this.on = ! this.on
            }
        })
    })
</script>
```
Given that we've registered the `darkMode` store and set `on` to "false", when the `<button>` is pressed, `on` will be "true" and the background color of the page will change to black.
## Single-value stores
If you don't need an entire object for a store, you can set and use any kind of data as a store.
Here's the example from above but using it more simply as a boolean value:
```alpine
<button x-data @click="$store.darkMode = ! $store.darkMode">Toggle Dark Mode</button>
...
<div x-data :class="$store.darkMode && 'bg-black'">
    ...
</div>
<script>
    document.addEventListener('alpine:init', () => {
        Alpine.store('darkMode', false)
    })
</script>
```
##### magics-watch
# $watch
You can "watch" a component property using the `$watch` magic method. For example:
```alpine
<div x-data="{ open: false }" x-init="$watch('open', value => console.log(value))">
    <button @click="open = ! open">Toggle Open</button>
</div>
```
In the above example, when the button is pressed and `open` is changed, the provided callback will fire and `console.log` the new value:
You can watch deeply nested properties using "dot" notation
```alpine
<div x-data="{ foo: { bar: 'baz' }}" x-init="$watch('foo.bar', value => console.log(value))">
    <button @click="foo.bar = 'bob'">Toggle Open</button>
</div>
```
When the `<button>` is pressed, `foo.bar` will be set to "bob", and "bob" will be logged to the console.
### Getting the "old" value
`$watch` keeps track of the previous value of the property being watched, You can access it using the optional second argument to the callback like so:
```alpine
<div x-data="{ open: false }" x-init="$watch('open', (value, oldValue) => console.log(value, oldValue))">
    <button @click="open = ! open">Toggle Open</button>
</div>
```
### Deep watching
`$watch` automatically watches from changes at any level but you should keep in mind that, when a change is detected, the watcher will return the value of the observed property, not the value of the subproperty that has changed.
```alpine
<div x-data="{ foo: { bar: 'baz' }}" x-init="$watch('foo', (value, oldValue) => console.log(value, oldValue))">
    <button @click="foo.bar = 'bob'">Update</button>
</div>
```
When the `<button>` is pressed, `foo.bar` will be set to "bob", and "{bar: 'bob'} {bar: 'baz'}" will be logged to the console (new and old value).
> ⚠️ Changing a property of a "watched" object as a side effect of the `$watch` callback will generate an infinite loop and eventually error. 
```alpine
<div x-data="{ foo: { bar: 'baz', bob: 'lob' }}" x-init="$watch('foo', value => foo.bob = foo.bar)">
    <button @click="foo.bar = 'bob'">Update</button>
</div>
```
##### plugins-anchor
# Anchor Plugin
Alpine's Anchor plugin allows you to easily anchor an element's positioning to another element on the page.
This functionality is useful when creating dropdown menus, popovers, dialogs, and tooltips with Alpine.
The "anchoring" functionality used in this plugin is provided by the [Floating UI](https://floating-ui.com/) project.
Initialize it within index.js (bundle):
```js
import Alpine from 'alpinejs'
import anchor from '@alpinejs/anchor'
Alpine.plugin(anchor)
...
```
## x-anchor
The primary API for using this plugin is the `x-anchor` directive.
To use this plugin, add the `x-anchor` directive to any element and pass it a reference to the element you want to anchor it's position to (often a button on the page).
By default, `x-anchor` will set the element's CSS to `position: absolute` and the appropriate `top` and `left` values. If the anchored element is normally displayed below the reference element but doesn't have room on the page, it's styling will be adjusted to render above the element.
For example, here's a simple dropdown anchored to the button that toggles it:
```alpine
<div x-data="{ open: false }">
    <button x-ref="button" @click="open = ! open">Toggle</button>
    <div x-show="open" x-anchor="$refs.button">
        Dropdown content
    </div>
</div>
```
<div x-data="{ open: false }" class="demo overflow-hidden">
    <div class="flex justify-center">
        <button x-ref="button" @click="open = ! open">Toggle</button>
    </div>
    <div x-show="open" x-anchor="$refs.button" class="bg-white rounded p-4 border shadow z-10">
        Dropdown content
    </div>
</div>
## Positioning
`x-anchor` allows you to customize the positioning of the anchored element using the following modifiers:
* Bottom: `.bottom`, `.bottom-start`, `.bottom-end`
* Top: `.top`, `.top-start`, `.top-end`
* Left: `.left`, `.left-start`, `.left-end`
* Right: `.right`, `.right-start`, `.right-end`
Here is an example of using `.bottom-start` to position a dropdown below and to the right of the reference element:
```alpine
<div x-data="{ open: false }">
    <button x-ref="button" @click="open = ! open">Toggle</button>
    <div x-show="open" x-anchor.bottom-start="$refs.button">
        Dropdown content
    </div>
</div>
```
<div x-data="{ open: false }" class="demo overflow-hidden">
    <div class="flex justify-center">
        <button x-ref="button" @click="open = ! open">Toggle</button>
    </div>
    <div x-show="open" x-anchor.bottom-start="$refs.button" class="bg-white rounded p-4 border shadow z-10">
        Dropdown content
    </div>
</div>
## Offset
You can add an offset to your anchored element using the `.offset.[px value]` modifier like so:
```alpine
<div x-data="{ open: false }">
    <button x-ref="button" @click="open = ! open">Toggle</button>
    <div x-show="open" x-anchor.offset.10="$refs.button">
        Dropdown content
    </div>
</div>
```
<div x-data="{ open: false }" class="demo overflow-hidden">
    <div class="flex justify-center">
        <button x-ref="button" @click="open = ! open">Toggle</button>
    </div>
    <div x-show="open" x-anchor.offset.10="$refs.button" class="bg-white rounded p-4 border shadow z-10">
        Dropdown content
    </div>
</div>
## Manual styling
By default, `x-anchor` applies the positioning styles to your element under the hood. If you'd prefer full control over styling, you can pass the `.no-style` modifer and use the `$anchor` magic to access the values inside another Alpine expression.
Below is an example of bypassing `x-anchor`'s internal styling and instead applying the styles yourself using `x-bind:style`:
```alpine
<div x-data="{ open: false }">
    <button x-ref="button" @click="open = ! open">Toggle</button>
    <div
        x-show="open"
        x-anchor.no-style="$refs.button"
        x-bind:style="{ position: 'absolute', top: $anchor.y+'px', left: $anchor.x+'px' }"
    >
        Dropdown content
    </div>
</div>
```
<div x-data="{ open: false }" class="demo overflow-hidden">
    <div class="flex justify-center">
        <button x-ref="button" @click="open = ! open">Toggle</button>
    </div>
    <div
        x-show="open"
        x-anchor.no-style="$refs.button"
        x-bind:style="{ position: 'absolute', top: $anchor.y+'px', left: $anchor.x+'px' }"
        class="bg-white rounded p-4 border shadow z-10"
    >
        Dropdown content
    </div>
</div>
## Anchor to an ID
The examples thus far have all been anchoring to other elements using Alpine refs.
Because `x-anchor` accepts a reference to any DOM element, you can use utilities like `document.getElementById()` to anchor to an element by its `id` attribute:
```alpine
<div x-data="{ open: false }">
    <button id="trigger" @click="open = ! open">Toggle</button>
    <div x-show="open" x-anchor="document.getElementById('trigger')">
        Dropdown content
    </div>
</div>
```
<div x-data="{ open: false }" class="demo overflow-hidden">
    <div class="flex justify-center">
        <button class="trigger" @click="open = ! open">Toggle</button>
    </div>
    <div x-show="open" x-anchor="document.querySelector('.trigger')">
        Dropdown content
    </div>
</div>
##### plugins-collapse
# Collapse Plugin
Alpine's Collapse plugin allows you to expand and collapse elements using smooth animations.
Because this behavior and implementation differs from Alpine's standard transition system, this functionality was made into a dedicated plugin.
Initialize it within index.js (bundle):
```js
import Alpine from 'alpinejs'
import collapse from '@alpinejs/collapse'
Alpine.plugin(collapse)
...
```
## x-collapse
The primary API for using this plugin is the `x-collapse` directive.
`x-collapse` can only exist on an element that already has an `x-show` directive. When added to an `x-show` element, `x-collapse` will smoothly "collapse" and "expand" the element when it's visibility is toggled by animating its height property.
For example:
```alpine
<div x-data="{ expanded: false }">
    <button @click="expanded = ! expanded">Toggle Content</button>
    <p x-show="expanded" x-collapse>
        ...
    </p>
</div>
```
<div x-data="{ expanded: false }" class="demo">
    <button @click="expanded = ! expanded">Toggle Content</button>
    <div x-show="expanded" x-collapse>
        <div class="pt-4">
            Reprehenderit eu excepteur ullamco esse cillum reprehenderit exercitation labore non. Dolore dolore ea dolore veniam sint in sint ex Lorem ipsum. Sint laborum deserunt deserunt amet voluptate cillum deserunt. Amet nisi pariatur sit ut id. Ipsum est minim est commodo id dolor sint id quis sint Lorem.
        </div>
    </div>
</div>
## Modifiers
### .duration
You can customize the duration of the collapse/expand transition by appending the `.duration` modifier like so:
```alpine
<div x-data="{ expanded: false }">
    <button @click="expanded = ! expanded">Toggle Content</button>
    <p x-show="expanded" x-collapse.duration.1000ms>
        ...
    </p>
</div>
```
<div x-data="{ expanded: false }" class="demo">
    <button @click="expanded = ! expanded">Toggle Content</button>
    <div x-show="expanded" x-collapse.duration.1000ms>
        <div class="pt-4">
            Reprehenderit eu excepteur ullamco esse cillum reprehenderit exercitation labore non. Dolore dolore ea dolore veniam sint in sint ex Lorem ipsum. Sint laborum deserunt deserunt amet voluptate cillum deserunt. Amet nisi pariatur sit ut id. Ipsum est minim est commodo id dolor sint id quis sint Lorem.
        </div>
    </div>
</div>
### .min
By default, `x-collapse`'s "collapsed" state sets the height of the element to `0px` and also sets `display: none;`.
Sometimes, it's helpful to "cut-off" an element rather than fully hide it. By using the `.min` modifier, you can set a minimum height for `x-collapse`'s "collapsed" state. For example:
```alpine
<div x-data="{ expanded: false }">
    <button @click="expanded = ! expanded">Toggle Content</button>
    <p x-show="expanded" x-collapse.min.50px>
        ...
    </p>
</div>
```
<div x-data="{ expanded: false }" class="demo">
    <button @click="expanded = ! expanded">Toggle Content</button>
    <div x-show="expanded" x-collapse.min.50px>
        <div class="pt-4">
            Reprehenderit eu excepteur ullamco esse cillum reprehenderit exercitation labore non. Dolore dolore ea dolore veniam sint in sint ex Lorem ipsum. Sint laborum deserunt deserunt amet voluptate cillum deserunt. Amet nisi pariatur sit ut id. Ipsum est minim est commodo id dolor sint id quis sint Lorem.
        </div>
    </div>
</div>
##### plugins-focus
> Notice: This Plugin was previously called "Trap". Trap's functionality has been absorbed into this plugin along with additional functionality. You can swap Trap for Focus without any breaking changes.
# Focus Plugin
Alpine's Focus plugin allows you to manage focus on a page.
> This plugin internally makes heavy use of the open source tool: [Tabbable](https://github.com/focus-trap/tabbable). Big thanks to that team for providing a much needed solution to this problem.
Initialize it in your bundle:
```js
import Alpine from 'alpinejs'
import focus from '@alpinejs/focus'
Alpine.plugin(focus)
...
```
## x-trap
Focus offers a dedicated API for trapping focus within an element: the `x-trap` directive.
`x-trap` accepts a JS expression. If the result of that expression is true, then the focus will be trapped inside that element until the expression becomes false, then at that point, focus will be returned to where it was previously.
For example:
```alpine
<div x-data="{ open: false }">
    <button @click="open = true">Open Dialog</button>
    <span x-show="open" x-trap="open">
        <p>...</p>
        <input type="text" placeholder="Some input...">
        <input type="text" placeholder="Some other input...">
        <button @click="open = false">Close Dialog</button>
    </span>
</div>
```
<div x-data="{ open: false }" class="demo">
    <div :class="open && 'opacity-50'">
        <button x-on:click="open = true">Open Dialog</button>
    </div>
    <div x-show="open" x-trap="open" class="mt-4 space-y-4 p-4 border bg-yellow-100" @keyup.escape.window="open = false">
        <strong>
            <div>Focus is now "trapped" inside this dialog, meaning you can only click/focus elements within this yellow dialog. If you press tab repeatedly, the focus will stay within this dialog.</div>
        </strong>
        <div>
            <input type="text" placeholder="Some input...">
        </div>
        <div>
            <input type="text" placeholder="Some other input...">
        </div>
        <div>
            <button @click="open = false">Close Dialog</button>
        </div>
    </div>
</div>
### Nesting dialogs
Sometimes you may want to nest one dialog inside another. `x-trap` makes this trivial and handles it automatically.
`x-trap` keeps track of newly "trapped" elements and stores the last actively focused element. Once the element is "untrapped" then the focus will be returned to where it was originally.
This mechanism is recursive, so you can trap focus within an already trapped element infinite times, then "untrap" each element successively.
Here is nesting in action:
```alpine
<div x-data="{ open: false }">
    <button @click="open = true">Open Dialog</button>
    <span x-show="open" x-trap="open">
        ...
        <div x-data="{ open: false }">
            <button @click="open = true">Open Nested Dialog</button>
            <span x-show="open" x-trap="open">
                ...
                <button @click="open = false">Close Nested Dialog</button>
            </span>
        </div>
        <button @click="open = false">Close Dialog</button>
    </span>
</div>
```
<div x-data="{ open: false }" class="demo">
    <div :class="open && 'opacity-50'">
        <button x-on:click="open = true">Open Dialog</button>
    </div>
    <div x-show="open" x-trap="open" class="mt-4 space-y-4 p-4 border bg-yellow-100" @keyup.escape.window="open = false">
        <div>
            <input type="text" placeholder="Some input...">
        </div>
        <div>
            <input type="text" placeholder="Some other input...">
        </div>
        <div x-data="{ open: false }">
            <div :class="open && 'opacity-50'">
                <button x-on:click="open = true">Open Nested Dialog</button>
            </div>
            <div x-show="open" x-trap="open" class="mt-4 space-y-4 p-4 border border-gray-500 bg-yellow-200" @keyup.escape.window="open = false">
                <strong>
                    <div>Focus is now "trapped" inside this nested dialog. You cannot focus anything inside the outer dialog while this is open. If you close this dialog, focus will be returned to the last known active element.</div>
                </strong>
                <div>
                    <input type="text" placeholder="Some input...">
                </div>
                <div>
                    <input type="text" placeholder="Some other input...">
                </div>
                <div>
                    <button @click="open = false">Close Nested Dialog</button>
                </div>
            </div>
        </div>
        <div>
            <button @click="open = false">Close Dialog</button>
        </div>
    </div>
</div>
### Modifiers
#### .inert
When building things like dialogs/modals, it's recommended to hide all the other elements on the page from screen readers when trapping focus.
By adding `.inert` to `x-trap`, when focus is trapped, all other elements on the page will receive `aria-hidden="true"` attributes, and when focus trapping is disabled, those attributes will also be removed.
```alpine
<body x-data="{ open: false }">
    <div x-trap.inert="open" ...>
        ...
    </div>
    <div>
        ...
    </div>
</body>
<body x-data="{ open: true }">
    <div x-trap.inert="open" ...>
        ...
    </div>
    <div aria-hidden="true">
        ...
    </div>
</body>
```
#### .noscroll
When building dialogs/modals with Alpine, it's recommended that you disable scrolling for the surrounding content when the dialog is open.
`x-trap` allows you to do this automatically with the `.noscroll` modifiers.
By adding `.noscroll`, Alpine will remove the scrollbar from the page and block users from scrolling down the page while a dialog is open.
For example:
```alpine
<div x-data="{ open: false }">
    <button @click="open = true">Open Dialog</button>
    <div x-show="open" x-trap.noscroll="open">
        Dialog Contents
        <button @click="open = false">Close Dialog</button>
    </div>
</div>
```
<div class="demo">
    <div x-data="{ open: false }">
        <button @click="open = true">Open Dialog</button>
        <div x-show="open" x-trap.noscroll="open" class="border mt-4 p-4">
            <div class="mb-4 text-bold">Dialog Contents</div>
            <p class="mb-4 text-gray-600 text-sm">Notice how you can no longer scroll on this page while this dialog is open.</p>
            <button class="mt-4" @click="open = false">Close Dialog</button>
        </div>
    </div>
</div>
#### .noreturn
Sometimes you may not want focus to be returned to where it was previously. Consider a dropdown that's triggered upon focusing an input, returning focus to the input on close will just trigger the dropdown to open again.
`x-trap` allows you to disable this behavior with the `.noreturn` modifier.
By adding `.noreturn`, Alpine will not return focus upon x-trap evaluating to false.
For example:
```alpine
<div x-data="{ open: false }" x-trap.noreturn="open">
    <input type="search" placeholder="search for something" />
    <div x-show="open">
        Search results
        <button @click="open = false">Close</button>
    </div>
</div>
```
<div class="demo">
    <div
        x-data="{ open: false }"
        x-trap.noreturn="open"
        @click.outside="open = false"
        @keyup.escape.prevent.stop="open = false"
    >
        <input type="search" placeholder="search for something"
            @focus="open = true"
            @keyup.escape.prevent="$el.blur()"
        />
        <div x-show="open">
            <div class="mb-4 text-bold">Search results</div>
            <p class="mb-4 text-gray-600 text-sm">Notice when closing this dropdown, focus is not returned to the input.</p>
            <button class="mt-4" @click="open = false">Close Dialog</button>
        </div>
    </div>
</div>
#### .noautofocus
By default, when `x-trap` traps focus within an element, it focuses the first focussable element within that element. This is a sensible default, however there are times where you may want to disable this behavior and not automatically focus any elements when `x-trap` engages.
By adding `.noautofocus`, Alpine will not automatically focus any elements when trapping focus.
## $focus
This plugin offers many smaller utilities for managing focus within a page. These utilities are exposed via the `$focus` magic.
| Property | Description |
| ---       | --- |
| `focus(el)`   | Focus the passed element (handling annoyances internally: using nextTick, etc.) |
| `focusable(el)`   | Detect whether or not an element is focusable |
| `focusables()`   | Get all "focusable" elements within the current element |
| `focused()`   | Get the currently focused element on the page |
| `lastFocused()`   | Get the last focused element on the page |
| `within(el)`   | Specify an element to scope the `$focus` magic to (the current element by default) |
| `first()`   | Focus the first focusable element |
| `last()`   | Focus the last focusable element |
| `next()`   | Focus the next focusable element |
| `previous()`   | Focus the previous focusable element |
| `noscroll()`   | Prevent scrolling to the element about to be focused |
| `wrap()`   | When retrieving "next" or "previous" use "wrap around" (ex. returning the first element if getting the "next" element of the last element) |
| `getFirst()`   | Retrieve the first focusable element |
| `getLast()`   | Retrieve the last focusable element |
| `getNext()`   | Retrieve the next focusable element |
| `getPrevious()`   | Retrieve the previous focusable element |
Let's walk through a few examples of these utilities in use. The example below allows the user to control focus within the group of buttons using the arrow keys. You can test this by clicking on a button, then using the arrow keys to move focus around:
```alpine
<div
    @keydown.right="$focus.next()"
    @keydown.left="$focus.previous()"
>
    <button>First</button>
    <button>Second</button>
    <button>Third</button>
</div>
```
<div class="demo">
<div
    x-data
    @keydown.right="$focus.next()"
    @keydown.left="$focus.previous()"
>
    <button class="focus:outline-none focus:ring-2 focus:ring-cyan-400">First</button>
    <button class="focus:outline-none focus:ring-2 focus:ring-cyan-400">Second</button>
    <button class="focus:outline-none focus:ring-2 focus:ring-cyan-400">Third</button>
</div>
(Click a button, then use the arrow keys to move left and right)
</div>
Notice how if the last button is focused, pressing "right arrow" won't do anything. Let's add the `.wrap()` method so that focus "wraps around":
```alpine
<div
    @keydown.right="$focus.wrap().next()"
    @keydown.left="$focus.wrap().previous()"
>
    <button>First</button>
    <button>Second</button>
    <button>Third</button>
</div>
```
<div class="demo">
<div
    x-data
    @keydown.right="$focus.wrap().next()"
    @keydown.left="$focus.wrap().previous()"
>
    <button class="focus:outline-none focus:ring-2 focus:ring-cyan-400">First</button>
    <button class="focus:outline-none focus:ring-2 focus:ring-cyan-400">Second</button>
    <button class="focus:outline-none focus:ring-2 focus:ring-cyan-400">Third</button>
</div>
(Click a button, then use the arrow keys to move left and right)
</div>
Now, let's add two buttons, one to focus the first element in the button group, and another focus the last element:
```alpine
<button @click="$focus.within($refs.buttons).first()">Focus "First"</button>
<button @click="$focus.within($refs.buttons).last()">Focus "Last"</button>
<div
    x-ref="buttons"
    @keydown.right="$focus.wrap().next()"
    @keydown.left="$focus.wrap().previous()"
>
    <button>First</button>
    <button>Second</button>
    <button>Third</button>
</div>
```
<div class="demo" x-data>
<button @click="$focus.within($refs.buttons).first()">Focus "First"</button>
<button @click="$focus.within($refs.buttons).last()">Focus "Last"</button>
<hr class="mt-2 mb-2"/>
<div
    x-ref="buttons"
    @keydown.right="$focus.wrap().next()"
    @keydown.left="$focus.wrap().previous()"
>
    <button class="focus:outline-none focus:ring-2 focus:ring-cyan-400">First</button>
    <button class="focus:outline-none focus:ring-2 focus:ring-cyan-400">Second</button>
    <button class="focus:outline-none focus:ring-2 focus:ring-cyan-400">Third</button>
</div>
</div>
Notice that we needed to add a `.within()` method for each button so that `$focus` knows to scope itself to a different element (the `div` wrapping the buttons).
##### plugins-intersect
# Intersect Plugin
Alpine's Intersect plugin is a convenience wrapper for [Intersection Observer](https://developer.mozilla.org/en-US/docs/Web/API/Intersection_Observer_API) that allows you to easily react when an element enters the viewport.
This is useful for: lazy loading images and other content, triggering animations, infinite scrolling, logging "views" of content, etc.
```js
import Alpine from 'alpinejs'
import intersect from '@alpinejs/intersect'
Alpine.plugin(intersect)
...
```
## x-intersect
The primary API for using this plugin is `x-intersect`. You can add `x-intersect` to any element within an Alpine component, and when that component enters the viewport (is scrolled into view), the provided expression will execute.
For example, in the following snippet, `shown` will remain `false` until the element is scrolled into view. At that point, the expression will execute and `shown` will become `true`:
```alpine
<div x-data="{ shown: false }" x-intersect="shown = true">
    <div x-show="shown" x-transition>
        I'm in the viewport!
    </div>
</div>
```
<div class="demo" style="height: 60px; overflow-y: scroll;" x-data x-ref="root">
    <a href="#" @click.prevent="$refs.root.scrollTo({ top: $refs.root.scrollHeight, behavior: 'smooth' })">Scroll Down 👇</a>
    <div style="height: 50vh"></div>
    <div x-data="{ shown: false }" x-intersect="shown = true" id="yoyo">
        <div x-show="shown" x-transition.duration.1000ms>
            I'm in the viewport!
        </div>
        <div x-show="! shown">&nbsp;</div>
    </div>
</div>
### x-intersect:enter
The `:enter` suffix is an alias of `x-intersect`, and works the same way:
```alpine
<div x-intersect:enter="shown = true">...</div>
```
You may choose to use this for clarity when also using the `:leave` suffix.
### x-intersect:leave
Appending `:leave` runs your expression when the element leaves the viewport.
```alpine
<div x-intersect:leave="shown = true">...</div>
```
> By default, this means the *whole element* is not in the viewport. Use `x-intersect:leave.full` to run your expression when only *parts of the element* are not in the viewport.
## Modifiers
### .once
Sometimes it's useful to evaluate an expression only the first time an element enters the viewport and not subsequent times. For example when triggering "enter" animations. In these cases, you can add the `.once` modifier to `x-intersect` to achieve this.
```alpine
<div x-intersect.once="shown = true">...</div>
```
### .half
Evaluates the expression once the intersection threshold exceeds `0.5`.
Useful for elements where it's important to show at least part of the element.
```alpine
<div x-intersect.half="shown = true">...</div> // when `0.5` of the element is in the viewport
```
### .full
Evaluates the expression once the intersection threshold exceeds `0.99`.
Useful for elements where it's important to show the whole element.
```alpine
<div x-intersect.full="shown = true">...</div> // when `0.99` of the element is in the viewport
```
### .threshold
Allows you to control the `threshold` property of the underlying `IntersectionObserver`:
This value should be in the range of "0-100". A value of "0" means: trigger an "intersection" if ANY part of the element enters the viewport (the default behavior). While a value of "100" means: don't trigger an "intersection" unless the entire element has entered the viewport.
Any value in between is a percentage of those two extremes.
For example if you want to trigger an intersection after half of the element has entered the page, you can use `.threshold.50`:
```alpine
<div x-intersect.threshold.50="shown = true">...</div> // when 50% of the element is in the viewport
```
If you wanted to trigger only when 5% of the element has entered the viewport, you could use: `.threshold.05`, and so on and so forth.
### .margin
Allows you to control the `rootMargin` property of the underlying `IntersectionObserver`.
This effectively tweaks the size of the viewport boundary. Positive values
expand the boundary beyond the viewport, and negative values shrink it inward. The values
work like CSS margin: one value for all sides; two values for top/bottom, left/right; or
four values for top, right, bottom, left. You can use `px` and `%` values, or use a bare number to
get a pixel value.
```alpine
<div x-intersect.margin.200px="loaded = true">...</div> // Load when the element is within 200px of the viewport
```
```alpine
<div x-intersect:leave.margin.10%.25px.25.25px="loaded = false">...</div> // Unload when the element gets within 10% of the top of the viewport, or within 25px of the other three edges
```
```alpine
<div x-intersect.margin.-100px="visible = true">...</div> // Mark as visible when element is more than 100 pixels into the viewport.
```
##### plugins-mask
# Mask Plugin
Alpine's Mask plugin allows you to automatically format a text input field as a user types.
This is useful for many different types of inputs: phone numbers, credit cards, dollar amounts, account numbers, dates, etc.
Initialize it within index.js (bundle):
```js
import Alpine from 'alpinejs'
import mask from '@alpinejs/mask'
Alpine.plugin(mask)
...
```
</div>
</div>
<button :aria-expanded="expanded" @click="expanded = ! expanded" class="text-cyan-600 font-medium underline">
    <span x-text="expanded ? 'Hide' : 'Show more'">Show</span> <span x-text="expanded ? '↑' : '↓'">↓</span>
</button>
</div>
## x-mask
The primary API for using this plugin is the `x-mask` directive.
Let's start by looking at the following simple example of a date field:
```alpine
<input x-mask="99/99/9999" placeholder="MM/DD/YYYY">
```
<div class="demo">
    <input x-data x-mask="99/99/9999" placeholder="MM/DD/YYYY">
</div>
Notice how the text you type into the input field must adhere to the format provided by `x-mask`. In addition to enforcing numeric characters, the forward slashes `/` are also automatically added if a user doesn't type them first.
The following wildcard characters are supported in masks:
| Wildcard | Description                      |
| -------- | -------------------------------- |
| `*`      | Any character                    |
| `a`      | Only alpha characters (a-z, A-Z) |
| `9`      | Only numeric characters (0-9)    |
## Dynamic Masks
Sometimes simple mask literals (i.e. `(999) 999-9999`) are not sufficient. In these cases, `x-mask:dynamic` allows you to dynamically generate masks on the fly based on user input.
Here's an example of a credit card input that needs to change it's mask based on if the number starts with the numbers "34" or "37" (which means it's an Amex card and therefore has a different format).
```alpine
<input x-mask:dynamic="
    $input.startsWith('34') || $input.startsWith('37')
        ? '9999 999999 99999' : '9999 9999 9999 9999'
">
```
As you can see in the above example, every time a user types in the input, that value is passed to the expression as `$input`. Based on the `$input`, a different mask is utilized in the field.
Try it for yourself by typing a number that starts with "34" and one that doesn't.
<div class="demo">
    <input x-data x-mask:dynamic="
        $input.startsWith('34') || $input.startsWith('37')
            ? '9999 999999 99999' : '9999 9999 9999 9999'
    ">
</div>
`x-mask:dynamic` also accepts a function as a result of the expression and will automatically pass it the `$input` as the first parameter. For example:
```alpine
<input x-mask:dynamic="creditCardMask">
<script>
function creditCardMask(input) {
    return input.startsWith('34') || input.startsWith('37')
        ? '9999 999999 99999'
        : '9999 9999 9999 9999'
}
</script>
```
## Money Inputs
Because writing your own dynamic mask expression for money inputs is fairly complex, Alpine offers a prebuilt one and makes it available as `$money()`.
Here is a fully functioning money input mask:
```alpine
<input x-mask:dynamic="$money($input)">
```
<div class="demo" x-data>
    <input type="text" x-mask:dynamic="$money($input)" placeholder="0.00">
</div>
If you wish to swap the periods for commas and vice versa (as is required in certain currencies), you can do so using the second optional parameter:
```alpine
<input x-mask:dynamic="$money($input, ',')">
```
<div class="demo" x-data>
    <input type="text" x-mask:dynamic="$money($input, ',')"  placeholder="0,00">
</div>
You may also choose to override the thousands separator by supplying a third optional argument:
```alpine
<input x-mask:dynamic="$money($input, '.', ' ')">
```
<div class="demo" x-data>
    <input type="text" x-mask:dynamic="$money($input, '.', ' ')"  placeholder="3 000.00">
</div>
You can also override the default precision of 2 digits by using any desired number of digits as the fourth optional argument:
```alpine
<input x-mask:dynamic="$money($input, '.', ',', 4)">
```
<div class="demo" x-data>
    <input type="text" x-mask:dynamic="$money($input, '.', ',', 4)"  placeholder="0.0001">
</div>
##### plugins-morph
# Morph Plugin
Alpine's Morph plugin allows you to "morph" an element on the page into the provided HTML template, all while preserving any browser or Alpine state within the "morphed" element.
This is useful for updating HTML from a server request without losing Alpine's on-page state. A utility like this is at the core of full-stack frameworks like [Laravel Livewire](https://laravel-livewire.com/) and [Phoenix LiveView](https://dockyard.com/blog/2018/12/12/phoenix-liveview-interactive-real-time-apps-no-need-to-write-javascript).
The best way to understand its purpose is with the following interactive visualization. Give it a try!
<div x-data="{ slide: 1 }" class="border rounded">
    <div>
        <img :src="'/img/morphs/morph'+slide+'.png'">
    </div>
    <div class="flex w-full justify-between" style="padding-bottom: 1rem">
        <div class="w-1/2 px-4">
            <button @click="slide = (slide === 1) ? 13 : slide - 1" class="w-full bg-cyan-400 rounded-full text-center py-3 font-bold text-white">Previous</button>
        </div>
        <div class="w-1/2 px-4">
            <button @click="slide = (slide % 13) + 1" class="w-full bg-cyan-400 rounded-full text-center py-3 font-bold text-white">Next</button>
        </div>
    </div>
</div>
```js
import Alpine from 'alpinejs'
import morph from '@alpinejs/morph'
window.Alpine = Alpine
Alpine.plugin(morph)
...
```
## Alpine.morph()
The `Alpine.morph(el, newHtml)` allows you to imperatively morph a dom node based on passed in HTML. It accepts the following parameters:
| Parameter | Description |
| ---       | --- |
| `el`      | A DOM element on the page. |
| `newHtml` | A string of HTML to use as the template to morph the dom element into. |
| `options` (optional) | An options object used mainly for [injecting lifecycle hooks](#lifecycle-hooks). |
Here's an example of using `Alpine.morph()` to update an Alpine component with new HTML: (In real apps, this new HTML would likely be coming from the server)
```alpine
<div x-data="{ message: 'Change me, then press the button!' }">
    <input type="text" x-model="message">
    <span x-text="message"></span>
</div>
<button>Run Morph</button>
<script>
    document.querySelector('button').addEventListener('click', () => {
        let el = document.querySelector('div')
        Alpine.morph(el, `
            <div x-data="{ message: 'Change me, then press the button!' }">
                <h2>See how new elements have been added</h2>
                <input type="text" x-model="message">
                <span x-text="message"></span>
                <h2>but the state of this component hasn't changed? Magical.</h2>
            </div>
        `)
    })
</script>
```
<div class="demo">
    <div x-data="{ message: 'Change me, then press the button!' }" id="morph-demo-1" class="space-y-2">
        <input type="text" x-model="message" class="w-full">
        <span x-text="message"></span>
    </div>
    <button id="morph-button-1" class="mt-4">Run Morph</button>
</div>
<script>
    document.querySelector('#morph-button-1').addEventListener('click', () => {
        let el = document.querySelector('#morph-demo-1')
        Alpine.morph(el, `
            <div x-data="{ message: 'Change me, then press the button!' }" id="morph-demo-1" class="space-y-2">
                <h4>See how new elements have been added</h4>
                <input type="text" x-model="message" class="w-full">
                <span x-text="message"></span>
                <h4>but the state of this component hasn't changed? Magical.</h4>
            </div>
        `)
    })
</script>
### Lifecycle Hooks
The "Morph" plugin works by comparing two DOM trees, the live element, and the passed in HTML.
Morph walks both trees simultaneously and compares each node and its children. If it finds differences, it "patches" (changes) the current DOM tree to match the passed in HTML's tree.
While the default algorithm is very capable, there are cases where you may want to hook into its lifecycle and observe or change its behavior as it's happening.
Before we jump into the available Lifecycle hooks themselves, let's first list out all the potential parameters they receive and explain what each one is:
| Parameter | Description |
| ---       | --- |
| `el` | This is always the actual, current, DOM element on the page that will be "patched" (changed by Morph). |
| `toEl` | This is a "template element". It's a temporary element representing what the live `el` will be patched to. It will never actually live on the page and should only be used for reference purposes. |
| `childrenOnly()` | This is a function that can be called inside the hook to tell Morph to skip the current element and only "patch" its children. |
| `skip()` | A function that when called within the hook will "skip" comparing/patching itself and the children of the current element. |
Here are the available lifecycle hooks (passed in as the third parameter to `Alpine.morph(..., options)`):
| Option | Description |
| ---       | --- |
| `updating(el, toEl, childrenOnly, skip)` | Called before patching the `el` with the comparison `toEl`.  |
| `updated(el, toEl)` | Called after Morph has patched `el`. |
| `removing(el, skip)` | Called before Morph removes an element from the live DOM. |
| `removed(el)` | Called after Morph has removed an element from the live DOM. |
| `adding(el, skip)` | Called before adding a new element. |
| `added(el)` | Called after adding a new element to the live DOM tree. |
| `key(el)` | A re-usable function to determine how Morph "keys" elements in the tree before comparing/patching. [More on that here](#keys) |
| `lookahead` | A boolean value telling Morph to enable an extra feature in its algorithm that "looks ahead" to make sure a DOM element that's about to be removed should instead just be "moved" to a later sibling. |
Here is code of all these lifecycle hooks for a more concrete reference:
```js
Alpine.morph(el, newHtml, {
    updating(el, toEl, childrenOnly, skip) {
        //
    },
    updated(el, toEl) {
        //
    },
    removing(el, skip) {
        //
    },
    removed(el) {
        //
    },
    adding(el, skip) {
        //
    },
    added(el) {
        //
    },
    key(el) {
        // By default Alpine uses the `key=""` HTML attribute.
        return el.id
    },
    lookahead: true, // Default: false
})
```
### Keys
Dom-diffing utilities like Morph try their best to accurately "morph" the original DOM into the new HTML. However, there are cases where it's impossible to determine if an element should be just changed, or replaced completely.
Because of this limitation, Morph has a "key" system that allows developers to "force" preserving certain elements rather than replacing them.
The most common use-case for them is a list of siblings within a loop. Below is an example of why keys are necessary sometimes:
```html
<ul>
    <li>Mark</li>
    <li>Tom</li>
    <li>Travis</li>
</ul>
<ul>
    <li>Travis</li>
    <li>Mark</li>
    <li>Tom</li>
</ul>
```
Given the above situation, Morph has no way to know that the "Travis" node has been moved in the DOM tree. It just thinks that "Mark" has been changed to "Travis" and "Travis" changed to "Tom".
This is not what we actually want, we want Morph to preserve the original elements and instead of changing them, MOVE them within the `<ul>`.
By adding keys to each node, we can accomplish this like so:
```html
<ul>
    <li key="1">Mark</li>
    <li key="2">Tom</li>
    <li key="3">Travis</li>
</ul>
<ul>
    <li key="3">Travis</li>
    <li key="1">Mark</li>
    <li key="2">Tom</li>
</ul>
```
Now that there are "keys" on the `<li>`s, Morph will match them in both trees and move them accordingly.
You can configure what Morph considers a "key" with the `key:` configuration option. [More on that here](#lifecycle-hooks)
##### plugins-persist
# Persist Plugin
Alpine's Persist plugin allows you to persist Alpine state across page loads.
This is useful for persisting search filters, active tabs, and other features where users will be frustrated if their configuration is reset after refreshing or leaving and revisiting a page.
Initialize it within index.js (bundle):
```js
import Alpine from 'alpinejs'
import persist from '@alpinejs/persist'
Alpine.plugin(persist)
...
```
## $persist
The primary API for using this plugin is the magic `$persist` method.
You can wrap any value inside `x-data` with `$persist` like below to persist its value across page loads:
```alpine
<div x-data="{ count: $persist(0) }">
    <button x-on:click="count++">Increment</button>
    <span x-text="count"></span>
</div>
```
<div class="demo">
    <div x-data="{ count: $persist(0) }">
        <button x-on:click="count++">Increment</button>
        <span x-text="count"></span>
    </div>
</div>
In the above example, because we wrapped `0` in `$persist()`, Alpine will now intercept changes made to `count` and persist them across page loads.
You can try this for yourself by incrementing the "count" in the above example, then refreshing this page and observing that the "count" maintains its state and isn't reset to "0".
## How does it work?
If a value is wrapped in `$persist`, on initialization Alpine will register its own watcher for that value. Now everytime that value changes for any reason, Alpine will store the new value in [localStorage](https://developer.mozilla.org/en-US/docs/Web/API/Window/localStorage).
Now when a page is reloaded, Alpine will check localStorage (using the name of the property as the key) for a value. If it finds one, it will set the property value from localStorage immediately.
You can observe this behavior by opening your browser devtool's localStorage viewer:
<a href="https://developer.chrome.com/docs/devtools/storage/localstorage/"><img src="/img/persist_devtools.png" alt="Chrome devtools showing the localStorage view with count set to 0"></a>
You'll observe that by simply visiting this page, Alpine already set the value of "count" in localStorage. You'll also notice it prefixes the property name "count" with "_x_" as a way of namespacing these values so Alpine doesn't conflict with other tools using localStorage.
Now change the "count" in the following example and observe the changes made by Alpine to localStorage:
```alpine
<div x-data="{ count: $persist(0) }">
    <button x-on:click="count++">Increment</button>
    <span x-text="count"></span>
</div>
```
<div class="demo">
    <div x-data="{ count: $persist(0) }">
        <button x-on:click="count++">Increment</button>
        <span x-text="count"></span>
    </div>
</div>
> `$persist` works with primitive values as well as with arrays and objects.
However, it is worth noting that localStorage must be cleared when the type of the variable changes.<br>
> Given the previous example, if we change count to a value of `$persist({ value: 0 })`, then localStorage must be cleared or the variable 'count' renamed.
## Setting a custom key
By default, Alpine uses the property key that `$persist(...)` is being assigned to ("count" in the above examples).
Consider the scenario where you have multiple Alpine components across pages or even on the same page that all use "count" as the property key.
Alpine will have no way of differentiating between these components.
In these cases, you can set your own custom key for any persisted value using the `.as` modifier like so:
```alpine
<div x-data="{ count: $persist(0).as('other-count') }">
    <button x-on:click="count++">Increment</button>
    <span x-text="count"></span>
</div>
```
Now Alpine will store and retrieve the above "count" value using the key "other-count".
Here's a view of Chrome Devtools to see for yourself:
<img src="/img/persist_custom_key_devtools.png" alt="Chrome devtools showing the localStorage view with count set to 0">
## Using a custom storage
By default, data is saved to localStorage, it does not have an expiration time and it's kept even when the page is closed.
Consider the scenario where you want to clear the data once the user close the tab. In this case you can persist data to sessionStorage using the `.using` modifier like so:
```alpine
<div x-data="{ count: $persist(0).using(sessionStorage) }">
    <button x-on:click="count++">Increment</button>
    <span x-text="count"></span>
</div>
```
You can also define your custom storage object exposing a getItem function and a setItem function. For example, you can decide to use a session cookie as storage doing so:
```alpine
<script>
    window.cookieStorage = {
        getItem(key) {
            let cookies = document.cookie.split(";");
            for (let i = 0; i < cookies.length; i++) {
                let cookie = cookies[i].split("=");
                if (key == cookie[0].trim()) {
                    return decodeURIComponent(cookie[1]);
                }
            }
            return null;
        },
        setItem(key, value) {
            document.cookie = key+' = '+encodeURIComponent(value)
        }
    }
</script>
<div x-data="{ count: $persist(0).using(cookieStorage) }">
    <button x-on:click="count++">Increment</button>
    <span x-text="count"></span>
</div>
```
## Using $persist with Alpine.data
If you want to use `$persist` with `Alpine.data`, you need to use a standard function instead of an arrow function so Alpine can bind a custom `this` context when it initially evaluates the component scope.
```js
Alpine.data('dropdown', function () {
    return {
        open: this.$persist(false)
    }
})
```
## Using the Alpine.$persist global
`Alpine.$persist` is exposed globally so it can be used outside of `x-data` contexts. This is useful to persist data from other sources such as `Alpine.store`.
```js
Alpine.store('darkMode', {
    on: Alpine.$persist(true).as('darkMode_on')
});
```
##### plugins-resize
# Resize Plugin
Alpine's Resize plugin is a convenience wrapper for the [Resize Observer](https://developer.mozilla.org/en-US/docs/Web/API/Resize_Observer_API) that allows you to easily react when an element changes size.
This is useful for: custom size-based animations, intelligent sticky positioning, conditionally adding attributes based on the element's size, etc.
Initialize it within index.js (bundle):
```js
import Alpine from 'alpinejs'
import resize from '@alpinejs/resize'
Alpine.plugin(resize)
...
```
## x-resize
The primary API for using this plugin is `x-resize`. You can add `x-resize` to any element within an Alpine component, and when that element is resized for any reason, the provided expression will execute with two magic properties: `$width` and `$height`.
For example, here's a simple example of using `x-resize` to display the width and height of an element as it changes size.
```alpine
<div
    x-data="{ width: 0, height: 0 }"
    x-resize="width = $width; height = $height"
>
    <p x-text="'Width: ' + width + 'px'"></p>
    <p x-text="'Height: ' + height + 'px'"></p>
</div>
```
<div class="demo">
    <div x-data="{ width: 0, height: 0 }" x-resize="width = $width; height = $height">
        <i>Resize your browser window to see the width and height values change.</i>
        <br><br>
        <p x-text="'Width: ' + width + 'px'"></p>
        <p x-text="'Height: ' + height + 'px'"></p>
    </div>
</div>
## Modifiers
### .document
It's often useful to observe the entire document's size, rather than a specific element. To do this, you can add the `.document` modifier to `x-resize`:
```alpine
<div x-resize.document="...">
```
<div class="demo">
    <div x-data="{ width: 0, height: 0 }" x-resize.document="width = $width; height = $height">
        <i>Resize your browser window to see the document width and height values change.</i>
        <br><br>
        <p x-text="'Width: ' + width + 'px'"></p>
        <p x-text="'Height: ' + height + 'px'"></p>
    </div>
</div>
##### plugins-sort
# Sort Plugin
Alpine's Sort plugin allows you to easily re-order elements by dragging them with your mouse.
This functionality is useful for things like Kanban boards, to-do lists, sortable table columns, etc.
The drag functionality used in this plugin is provided by the [SortableJS](https://github.com/SortableJS/Sortable) project.
Initialize it within index.js (bundle):
```js
import Alpine from 'alpinejs'
import sort from '@alpinejs/sort'
Alpine.plugin(sort)
...
```
## Basic usage
The primary API for using this plugin is the `x-sort` directive. By adding `x-sort` to an element, its children containing `x-sort:item` become sortable—meaning you can drag them around with your mouse, and they will change positions.
```alpine
<ul x-sort>
    <li x-sort:item>foo</li>
    <li x-sort:item>bar</li>
    <li x-sort:item>baz</li>
</ul>
```
<div x-data>
    <ul x-sort>
        <li x-sort:item class="cursor-pointer">foo</li>
        <li x-sort:item class="cursor-pointer">bar</li>
        <li x-sort:item class="cursor-pointer">baz</li>
    </ul>
</div>
## Sort handlers
You can react to sorting changes by passing a handler function to `x-sort` and adding keys to each item using `x-sort:item`. Here is an example of a simple handler function that shows an alert dialog with the changed item's key and its new position:
```alpine
<ul x-sort="alert($item + ' - ' + $position)">
    <li x-sort:item="1">foo</li>
    <li x-sort:item="2">bar</li>
    <li x-sort:item="3">baz</li>
</ul>
```
<div x-data>
    <ul x-sort="alert($item + ' - ' + $position)">
        <li x-sort:item="1" class="cursor-pointer">foo</li>
        <li x-sort:item="2" class="cursor-pointer">bar</li>
        <li x-sort:item="3" class="cursor-pointer">baz</li>
    </ul>
</div>
The `x-sort` handler will be called every time the sort order of the items change. The `$item` magic will contain the key of the sorted element (derived from `x-sort:item`), and `$position` will contain the new position of the item (starting at index `0`).
You can also pass a handler function to `x-sort` and that function will receive the `item` and `position` as the first and second parameter:
```alpine
<div x-data="{ handle: (item, position) => { ... } }">
    <ul x-sort="handle">
        <li x-sort:item="1">foo</li>
        <li x-sort:item="2">bar</li>
        <li x-sort:item="3">baz</li>
    </ul>
</div>
```
Handler functions are often used to persist the new order of items in the database so that the sorting order of a list is preserved between page refreshes.
## Sorting groups
This plugin allows you to drag items from one `x-sort` sortable list into another one by adding a matching `x-sort:group` value to both lists:
```alpine
<div>
    <ul x-sort x-sort:group="todos">
        <li x-sort:item="1">foo</li>
        <li x-sort:item="2">bar</li>
        <li x-sort:item="3">baz</li>
    </ul>
    <ol x-sort x-sort:group="todos">
        <li x-sort:item="4">foo</li>
        <li x-sort:item="5">bar</li>
        <li x-sort:item="6">baz</li>
    </ol>
</div>
```
Because both sortable lists above use the same group name (`todos`), you can drag items from one list onto another.
> When using sort handlers like `x-sort="handle"` and dragging an item from one group to another, only the destination list's handler will be called with the key and new position.
## Drag handles
By default, each `x-sort:item` element is draggable by clicking and dragging anywhere within it. However, you may want to designate a smaller, more specific element as the "drag handle" so that the rest of the element can be interacted with like normal, and only the handle will respond to mouse dragging:
```alpine
<ul x-sort>
    <li x-sort:item>
        <span x-sort:handle> - </span>foo
    </li>
    <li x-sort:item>
        <span x-sort:handle> - </span>bar
    </li>
    <li x-sort:item>
        <span x-sort:handle> - </span>baz
    </li>
</ul>
```
<div x-data>
    <ul x-sort>
        <li x-sort:item>
            <span x-sort:handle class="cursor-pointer"> - </span>foo
        </li>
        <li x-sort:item>
            <span x-sort:handle class="cursor-pointer"> - </span>bar
        </li>
        <li x-sort:item>
            <span x-sort:handle class="cursor-pointer"> - </span>baz
        </li>
    </ul>
</div>
As you can see in the above example, the hyphen "-" is draggable, but the item text ("foo") is not.
## Ghost elements
When a user drags an item, the element will follow their mouse to appear as though they are physically dragging the element.
By default, a "hole" (empty space) will be left in the original element's place during the drag.
If you would like to show a "ghost" of the original element in its place instead of an empty space, you can add the `.ghost` modifier to `x-sort`:
```alpine
<ul x-sort.ghost>
    <li x-sort:item>foo</li>
    <li x-sort:item>bar</li>
    <li x-sort:item>baz</li>
</ul>
```
<div x-data>
    <ul x-sort.ghost>
        <li x-sort:item class="cursor-pointer">foo</li>
        <li x-sort:item class="cursor-pointer">bar</li>
        <li x-sort:item class="cursor-pointer">baz</li>
    </ul>
</div>
### Styling the ghost element
By default, the "ghost" element has a `.sortable-ghost` CSS class attached to it while the original element is being dragged.
This makes it easy to add any custom styling you would like:
```alpine
<style>
.sortable-ghost {
    opacity: .5 !important;
}
</style>
<ul x-sort.ghost>
    <li x-sort:item>foo</li>
    <li x-sort:item>bar</li>
    <li x-sort:item>baz</li>
</ul>
```
<div x-data>
    <ul x-sort.ghost x-sort:config="{ ghostClass: 'opacity-50' }">
        <li x-sort:item class="cursor-pointer">foo</li>
        <li x-sort:item class="cursor-pointer">bar</li>
        <li x-sort:item class="cursor-pointer">baz</li>
    </ul>
</div>
## Sorting class on body
While an element is being dragged around, Alpine will automatically add a `.sorting` class to the `<body>` element of the page.
This is useful for styling any element on the page conditionally using only CSS.
For example you could have a warning that only displays while a user is sorting items:
```html
<div id="sort-warning">
    Page functionality is limited while sorting
</div>
```
To show this only while sorting, you can use the `body.sorting` CSS selector:
```css
#sort-warning {
    display: none;
}
body.sorting #sort-warning {
    display: block;
}
```
## CSS hover bug
Currently, there is a bug in Chrome and Safari (not Firefox) that causes issues with hover styles.
Consider HTML like the following, where each item in the list is styled differently based on a hover state (here we're using Tailwind's `.hover` class to conditionally add a border):
```html
<div x-sort>
    <div x-sort:item class="hover:border">foo</div>
    <div x-sort:item class="hover:border">bar</div>
    <div x-sort:item class="hover:border">baz</div>
</div>
```
If you drag one of the elements in the list below you will see that the hover effect will be errantly applied to any element in the original element's place:
<div x-data>
    <ul x-sort class="flex flex-col items-start">
        <li x-sort:item class="hover:border border-black cursor-pointer">foo</li>
        <li x-sort:item class="hover:border border-black cursor-pointer">bar</li>
        <li x-sort:item class="hover:border border-black cursor-pointer">baz</li>
    </ul>
</div>
To fix this, you can leverage the `.sorting` class applied to the body while sorting to limit the hover effect to only be applied while `.sorting` does NOT exist on `body`.
Here is how you can do this directly inline using Tailwind arbitrary variants:
```html
<div x-sort>
    <div x-sort:item class="[body:not(.sorting)_&]:hover:border">foo</div>
    <div x-sort:item class="[body:not(.sorting)_&]:hover:border">bar</div>
    <div x-sort:item class="[body:not(.sorting)_&]:hover:border">baz</div>
</div>
```
Now you can see below that the hover effect is only applied to the dragging element and not the others in the list.
<div x-data>
    <ul x-sort class="flex flex-col items-start">
        <li x-sort:item class="[body:not(.sorting)_&]:hover:border border-black cursor-pointer">foo</li>
        <li x-sort:item class="[body:not(.sorting)_&]:hover:border border-black cursor-pointer">bar</li>
        <li x-sort:item class="[body:not(.sorting)_&]:hover:border border-black cursor-pointer">baz</li>
    </ul>
</div>
## Custom configuration
Alpine chooses sensible defaults for configuring SortableJS under the hood. However, you can add or override any of these options yourself using `x-sort:config`:
```alpine
<ul x-sort x-sort:config="{ animation: 0 }">
    <li x-sort:item>foo</li>
    <li x-sort:item>bar</li>
    <li x-sort:item>baz</li>
</ul>
```
<div x-data>
    <ul x-sort x-sort:config="{ animation: 0 }">
        <li x-sort:item class="cursor-pointer">foo</li>
        <li x-sort:item class="cursor-pointer">bar</li>
        <li x-sort:item class="cursor-pointer">baz</li>
    </ul>
</div>
> Any config options passed will overwrite Alpine defaults. In this case of `animation`, this is fine, however be aware that overwriting `handle`, `group`, `filter`, `onSort`, `onStart`, or `onEnd` may break functionality.
