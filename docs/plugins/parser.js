export function parseCodeText(codeText) {
    const language = detectLanguage(codeText);

    switch (language) {
        case "javascript":
            return {
                language: "javascript",
                functions: parseJavaScriptFunctions(codeText),
            };
        case "python":
            return {
                language: "python",
                functions: parsePythonFunctions(codeText),
            };
        case "go":
            return { language: "go", functions: parseGoFunctions(codeText) };
        case "rust":
            return {
                language: "rust",
                functions: parseRustFunctions(codeText),
            };
        default:
            return { language: "unknown", functions: [] };
    }
}

export function detectLanguage(codeText) {
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

export function detectFileType(inputString) {
    // Split the input string by line breaks
    const lines = inputString.split("\n");
    // Iterate over each line
    for (const line of lines) {
        // Extract the file path from the line
        const filePathMatch = line.match(/\/[^:]+/);
        if (filePathMatch) {
            const filePath = filePathMatch[0];
            // Extract the file extension
            const extension = filePath.split(".").pop();
            switch (extension) {
                case "go":
                    return "Go";
                case "rs":
                    return "Rust";
                case "py":
                    return "Python";
                case "js":
                    return "JavaScript";
                default:
                    // If no specific file type matches, continue to the next line
                    continue;
            }
        }
    }
    // If no file type is detected, return 'Unknown'
    return "Unknown";
}

export function parseError(errorMessage) {
    const lines = errorMessage.split("\n");
    const parsedStack = [];

    // Regular expression pattern to match the file path and line number
    const filePathRegex = /([^ ]+):(\d+)/;

    let functionName = null;

    for (let i = 0; i < lines.length; i++) {
        const line = lines[i];

        // Attempt to match the regular expression to extract file path and line number
        const filePathMatch = line.match(filePathRegex);
        if (filePathMatch) {
            const [, path, lineNum] = filePathMatch;

            // If we found a file path and line number, check the line above to extract function name
            if (i > 0 && functionName) {
                const name = functionName.split("(")[0].trim();
                parsedStack.push({
                    line: functionName.trim(),
                    name,
                    path,
                    lineNum: parseInt(lineNum),
                });
                functionName = null;
            }
        } else {
            // If the current line doesn't match file path pattern, check if it contains a function name
            functionName = line.trim();
        }
    }

    return parsedStack.reverse();
}

export function parseGoFunctions(code) {
    const functionRegex =
        /func\s*(?:\(([^)]*)\))?\s*(\w+)\(([^)]*)\)([^{]*)\{/g;
    let functions = [];
    let match;

    while ((match = functionRegex.exec(code)) !== null) {
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

export function parseJavaScriptFunctions(code) {
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

export function parsePythonFunctions(code) {
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

export function parseRustFunctions(code) {
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
