function parseCodeText(codeText) {
  const language = detectLanguage(codeText);

  switch (language) {
    case "javascript":
      return {
        language: "javascript",
        functions: parseJavaScriptFunctions(codeText),
      };
    case "python":
      return { language: "python", functions: parsePythonFunctions(codeText) };
    case "go":
      return { language: "go", functions: parseGoFunctions(codeText) };
    case "rust":
      return { language: "rust", functions: parseRustFunctions(codeText) };
    default:
      return { language: "unknown", functions: [] };
  }
}

function detectLanguage(codeText) {
  // For simplicity, let's assume we have some predefined patterns for each language
  if (/function/.test(codeText)) {
    return "javascript";
  } else if (/def\s+\w+\(.*\):/.test(codeText)) {
    return "python";
  } else if (/func\s+\w+\(.*\)/.test(codeText)) {
    return "go";
  } else if (/fn\s+\w+\(.*\)/.test(codeText)) {
    return "rust";
  } else {
    return "unknown";
  }
}
function parseGoFunctions(code) {
  const functionRegex = /func\s*(?:\(([^)]*)\))?\s*(\w+)\(([^)]*)\)([^{]*)\{/g;
  let functions = [];
  let match;

  while ((match = functionRegex.exec(code)) !== null) {
    console.log(match);
    let name = match[2];
    let parameters = match[3].split(",").map((param) => param.trim());
    if (match[1]) {
      parameters.unshift(match[1]);
    }
    for (let i = 0; i < parameters.length; i++) {
      const element = parameters[i].split(" ");
      parameters[i] = {
        name: element[0],
        type: element[1],
      };
    }
    let ret = match[4].trim();
    functions.push({ name, parameters, return: ret, definition: match[0] });
  }

  return functions;
}

function parseJavaScriptFunctions(code) {
  const functionRegex = /(?:function\s+)?(\w+)\s*\((.*?)\)\s*{/g;
  let functions = [];
  let match;

  while ((match = functionRegex.exec(code)) !== null) {
    let name = match[1];
    let parameters = match[2].split(",").map((param) => param.trim());
    functions.push({ name, parameters, definition: match[0] });
  }

  return functions;
}

function parsePythonFunctions(code) {
  const functionRegex = /def\s+(\w+)\((.*?)\):/g;
  let functions = [];
  let match;

  while ((match = functionRegex.exec(code)) !== null) {
    let name = match[1];
    let parameters = match[2].split(",").map((param) => param.trim());
    functions.push({ name, parameters, definition: match[0] });
  }

  return functions;
}

function parseRustFunctions(code) {
  const functionRegex = /fn\s+(\w+)\((.*?)\)\s*(?:->\s*(\w+))?/g;
  let functions = [];
  let match;

  while ((match = functionRegex.exec(code)) !== null) {
    let name = match[1];
    let parameters = match[2].split(",").map((param) => {
      const parts = param.trim().split(":");
      return { name: parts[0], type: parts[1] ? parts[1].trim() : null };
    });
    let returnType = match[3] ? match[3].trim() : null;
    functions.push({
      name,
      parameters,
      return: returnType,
      definition: match[0],
    });
  }

  return functions;
}

(() => {
  let code = document.querySelectorAll("pre code");

  let functions = {};

  for (let i = 0; i < code.length; i++) {
    let res = parseCodeText(code[i].innerText);
    if (functions[res.language] == undefined) {
      functions[res.language] = [];
    }
    functions[res.language].push(...res.functions);
  }

  let hs = document.querySelectorAll("h1,h2,h3,h4,h5,h6");

  let tagMatch = /\w+\?\((\w+)\)/;

  for (let i = 0; i < hs.length; i++) {
    const element = hs[i];
    if (tagMatch.test(element.innerText)) {
      let match = tagMatch.exec(element.innerText);
      let name = match[0].slice(0, match[0].indexOf("?"));
      let lang = match[1];
      element.innerHTML = element.innerHTML.slice(
        0,
        element.innerHTML.indexOf("?")
      );

      let func = functions[lang].filter((e) => {
        return e.name == name;
      })[0];

      // Inject blockquote

      let bq = document.createElement("blockquote");
      let p = document.createElement("p");
      p.innerText = func.definition;
      bq.appendChild(p);
      element.parentElement.insertBefore(bq, element);

      let table = document.createElement("table");
      let head = document.createElement("thead");

      let headtr = document.createElement("tr");
      let th1 = document.createElement("th");
      th1.innerHTML = "Name";
      headtr.appendChild(th1);
      let th2 = document.createElement("th");
      th2.innerHTML = "Type";
      headtr.appendChild(th2);

      head.appendChild(headtr);

      table.appendChild(head);

      let body = document.createElement("tbody");

      for (let i = 0; i < func.parameters.length; i++) {
        const element = func.parameters[i];
        let bodytr = document.createElement("tr");
        let th1 = document.createElement("td");
        th1.innerHTML = element.name;
        bodytr.appendChild(th1);
        let th2 = document.createElement("td");
        th2.innerHTML = element.type;
        bodytr.appendChild(th2);

        body.appendChild(bodytr);
      }

      let bodytr = document.createElement("tr");
      th1 = document.createElement("th");
      th1.innerHTML = "return";
      bodytr.appendChild(th1);
      th2 = document.createElement("th");
      th2.innerHTML = func.return;
      bodytr.appendChild(th2);

      body.appendChild(bodytr);

      table.appendChild(body);

      element.insertAdjacentElement("afterend", table);
    }
  }
})();
