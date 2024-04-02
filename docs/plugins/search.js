/**
 * It fetches the HTML of the DuckDuckGo search page, parses it, and returns the results
 * @param query - The query to search for
 * @returns An array of objects with the following properties:
 *     title: The title of the result
 *     description: The description of the result
 *     url: The url of the result
 */
export async function search(query) {
    const html = await fetch(
        `https://cors.lowsh.workers.dev/?https://lite.duckduckgo.com/lite/?q=${encodeURIComponent(
            query
        )}`
    );
    let text = await html.text();
    let doc = parseHTML(text);
    let sponsored = [
        ...doc.querySelectorAll("tr[class='result-sponsored']"),
    ].pop();
    let trs = [...doc.querySelectorAll("tr")];
    let rawRes = [...chunks(trs.slice(trs.indexOf(sponsored) + 1), 4)];

    let results = [];
    for (let i = 0; i < rawRes.length; i++) {
        const group = rawRes[i];
        if (group.length == 4) {
            results.push({
                title: group[1].querySelector("a").textContent,
                description: group[2].querySelector(
                    "td[class='result-snippet']"
                ).textContent,
                url:
                    "http://" +
                    group[3].querySelector("span[class='link-text']")
                        .textContent,
            });
        }
    }
    return results;
}

function* chunks(arr, n) {
    for (let i = 0; i < arr.length; i += n) {
        yield arr.slice(i, i + n);
    }
}

function parseHTML(html) {
    const template = document.createElement("template");
    template.innerHTML = html;
    return template.content;
}
