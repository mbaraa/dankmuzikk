htmx.defineExtension("json-enc",{onEvent:function(e,n){"htmx:configRequest"===e&&(n.detail.headers["Content-Type"]="application/json")},encodeParameters:function(e,n,t){for(let o of(n={},t))if(console.log("type",o.type,"name",o.name,"value",o.value),o.name&&o.value){let a=!1;"number"===o.type&&(a=!0),n[o.name]=a?Number(o.value):o.value}return console.log(n),e.overrideMimeType("text/json"),JSON.stringify(n)}});
