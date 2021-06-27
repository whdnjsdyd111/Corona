let cag;

const renderTables = async () => {
    cag = await cagItems();

    let content = ""
    for (let i = 0; i < cag.length; i++) {
        content += makeTable(cag, i);
    }

    $('#tables').html(content);
}

function makeTable(arr, i) {
    let half = arr.length / 2;
    return `<tr>
                <td>${arr[i].Gubun}</td>
                <td>${arr[i].ConfCase}</td>
                <td>${arr[i].ConfCaseRate}</td>
                <td>${arr[i].Death}</td>
                <td>${arr[i].DeathRate}</td>
            </tr>`;
}

$(function() {
    renderTables();

    $('#page_all').click(async () => {
        await movePage("index")
    });

    $('#page_map').click( async () => {
        await movePage("map")
    });
})