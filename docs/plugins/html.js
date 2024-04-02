export let style = function () {
    let css = templateToString(...arguments);
    let newCSS = renameCSSClasses(css);
    let s = document.createElement("style");
    s.innerHTML = newCSS.modifiedCSS;
    document.head.appendChild(s);
    return newCSS.classMap;
};

export function inline(el, styles) {
    for (const key in styles) {
        el.style[key] = styles[key];
    }
}

export let Fmt = function (strings, ...values) {
    const elements = strings.map((e, i) => {
        return { el: values[i], depth: e.length };
    });

    const stack = [elements[0]];

    for (let i = 1; i < elements.length; i++) {
        const { el, depth } = elements[i];

        if (el == undefined) {
            break;
        }

        while (stack.length > 1 && depth <= stack[stack.length - 1].depth) {
            stack.pop();
        }

        stack[stack.length - 1].el.appendChild(el);
        stack.push({ depth, el });
    }

    return values[0];
};

export function DOM(type) {
    return function () {
        let el = document.createElement(type);
        let props = parseHTMLProperties(...arguments);
        for (let prop of Object.keys(props)) {
            if (el[prop] !== undefined) {
                el[prop] = props[prop];
            } else {
                el.setAttribute(prop, props[prop]);
            }
        }
        el.on = (evt, handler) => {
            el.addEventListener(evt, handler);
            return el;
        };
        el.bind = (state, name) => {
            if (typeof el.value != "undefined") {
                el.addEventListener("input", () => {
                    state.val(name, el.value);
                });
                state.f((d) => {
                    el.value = d[name];
                });
            } else {
                state.f((d) => {
                    el.innerHTML = d[name];
                });
            }

            return el;
        };
        el.add = el.appendChild;
        el.clear = () => {
            el.innerHTML = "";
        };
        return el;
    };
}
function parseHTMLProperties(template, ...values) {
    const properties = {};
    const templateString = template.reduce((acc, part, i) => {
        const value = values[i];
        if (typeof value === "string") {
            acc += part + value;
        } else {
            // If the value is not a string, treat it as a property name
            const propName = value;
            if (propName !== undefined) {
                acc += part + `${propName}`;
            } else {
                acc += part;
            }
        }
        return acc;
    }, "");
    // Regular expression to match attribute="value" or attribute='value' or attribute
    const attributeRegex = /([\w-]+)(?:=(?:"([^"]*)"|'([^']*)'))?/g;
    templateString.replace(
        attributeRegex,
        (match, name, doubleQuotedValue, singleQuotedValue) => {
            const value = doubleQuotedValue || singleQuotedValue || true;

            properties[name] = value;
        }
    );

    return properties;
}

function templateToString(strings, ...values) {
    let result = "";
    for (let i = 0; i < strings.length; i++) {
        result += strings[i];
        if (i < values.length) {
            result += values[i];
        }
    }
    return result;
}

function generateRandomHash(length) {
    let result = "";
    const characters =
        "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    const charactersLength = characters.length;
    for (let i = 0; i < length; i++) {
        result += characters.charAt(
            Math.floor(Math.random() * charactersLength)
        );
    }
    return result;
}

function renameCSSClasses(css) {
    const classMap = {};
    const cssWithModifiedClasses = css.replace(
        /(\.)([\w-]+(?![^{}]*\}))/g,
        (match, selectorType, selectorName) => {
            if (selectorType === "." || selectorType === "#") {
                if (!classMap[selectorName]) {
                    classMap[
                        selectorName
                    ] = `${selectorName}-${generateRandomHash(8)}`;
                }
                return `${selectorType}${classMap[selectorName]}`;
            } else {
                return match; // Preserve properties like background
            }
        }
    );

    return {
        modifiedCSS: cssWithModifiedClasses,
        classMap,
    };
}
