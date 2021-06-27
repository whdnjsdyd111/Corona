const allSido = {
    "경기도": ["서울", "경기", "인천"],
    "강원도": ["강원"],
    "충청도": ["충남", "충북", "세종", "대전"],
    "전라도": ["전북", "전남", "광주"],
    "경상도": ["경북", "대구", "경남", "울산", "부산"],
    "제주도": ["제주"],
    "검역": ["검역"],
}

let defRate;

const renderStatistic = async () => {
    defRate = await accDefRate();
    setPatient(defRate);
}

let titles;

const renderTitles = async () => {
    titles = await initTitles();
    setTitles(titles);
}

let arr;

const renderTables = async () => {
    arr = await stateItems();

    let content = ""
    for (let i = 0; i < arr.length / 2; i++) {
        content += makeTable(arr, i);
    }

    setCure(Math.floor((arr[18].IsolClearCnt / arr[18].DefCnt) * 10000) / 100);
    $('#tables').html(content);
}

let selectedDo = "전체보기";

$(function() {
    renderStatistic()
    renderTitles()
    renderTables()

    $(document).on('change', '#do', function() {
        selectedDo = $(this).val();
        setAllTable(selectedDo);
    });

    $(document).on('change', '#sido', function() {
        let item = $(this).val();

        if (item === "전체보기") {
            setAllTable(selectedDo);
            return;
        }

        let sido;
        arr.some((v, i) => {
            if (item === v.Gubun) {
                $('#tables').html(makeTable(arr, i));
                sido = i;
                return true;
            }
        });
        let sidoTitle = [
            [0, 0],
            [0, 0],
            [0, 0]
        ];
        sidoTitle[0][0] = arr[sido].DefCnt;
        sidoTitle[1][0] = arr[sido].IsolClearCnt;
        sidoTitle[2][0] = arr[sido].DeathCnt;
        sidoTitle[0][1] = sidoTitle[0][0] - arr[sido + arr.length / 2].DefCnt;
        sidoTitle[1][1] = sidoTitle[1][0] - arr[sido + arr.length / 2].IsolClearCnt;
        sidoTitle[2][1] = sidoTitle[2][0] - arr[sido + arr.length / 2].DeathCnt;
        setTitles(sidoTitle)
        setCure(Math.floor(sidoTitle[1][0] / sidoTitle[0][0] * 10000) / 100);
    });

    $('#page_map').click( async () => {
        await movePage("map")
    });

    $('#page_ageGender').click(async () => {
        await movePage("ageGender")
    });
})

function makeComma(str) {
    str = String(str);

    return str.replace(/(\d)(?=(?:\d{3})+(?!\d))/g, '$1,');
}

function upOrDown(today, yesterday) {
    let check = today - yesterday;
    let up = check > 0;
    if (check === 0) return '';
    return `<i style='--color: var(${up ? "--i_red" : "--i_blue"})'>${up ? '↑' : '↓'} ${Math.abs(today - yesterday)}</i>`;
}

function makeTable(arr, i) {
    let half = arr.length / 2;
    return `<tr>
                <td>${arr[i].Gubun}</td>
                <td>${arr[i].DefCnt}${upOrDown(arr[i].IncDec, arr[i + half].IncDec)}</td>
                <td>${arr[i].OverFlowCnt}${upOrDown(arr[i].OverFlowCnt, arr[i + half].OverFlowCnt)}</td>
                <td>${arr[i].DeathCnt}${upOrDown(arr[i].DeathCnt, arr[i + half].DeathCnt)}</td>
                <td>${arr[i].IsolIngCnt}${upOrDown(arr[i].IsolIngCnt, arr[i + half].IsolIngCnt)}</td>
                <td>${Math.floor(arr[i].QurRate * 100) / 100}</td>
            </tr>`;
}

function setTitles(titles) {
    for (let i = 0; i < 3; i++) {
        let title = $($('.title')[i]);
        title.children('p').text(makeComma(titles[i][0]));
        let tagI = title.children('i');
        if (titles[i][1] > 0) {
            tagI.text(`↑ ${makeComma(titles[i][1])}`);
            tagI.get(0).style.setProperty('--color', 'var(--i_red)');
        } else {
            tagI.text(`↓ ${makeComma(titles[i][1])}`);
            tagI.get(0).style.setProperty('--color', 'var(--i_blue)');
        }
    }
}

function setPatient(rate) {
    $('#defRateCir').removeClass('defCir');
    $('#defRate').text(`${rate}%`);
    $('#defRateCir').get(0).style.setProperty("--i", rate);
    setTimeout(function() {
        $('#defRateCir').addClass('defCir');
    }, 10);
}

function setCure(rate) {
    $('#cureRateCir').removeClass('cureCir');
    $('#cureRate').text(`${rate}%`);
    $('#cureRateCir').get(0).style.setProperty("--i", rate);
    setTimeout(function() {
        $('#cureRateCir').addClass('cureCir');
    }, 10);
}

function setAllTable(item) {
    let sido = $('#sido');
    sido.html('<option selected>전체보기</option>');

    let content = "";
    if (item === "전체보기") {
        for (let i = 0; i < arr.length / 2; i++) {
            content += makeTable(arr, i);
        }
        $('#tables').html(content);
        setTitles(titles);
        setPatient(defRate);
        setCure(Math.floor(arr[18].IsolClearCnt / arr[18].DefCnt * 10000) / 100);
        return;
    }

    let sidoTitles = [
        [0, 0], // 확진자
        [0, 0],
        [0, 0]
    ];
    for (let i = 0; i < allSido[item].length; i++) {
        let index;
        arr.some((v, ix) => {
            if (v.Gubun === allSido[item][i]) {
                index = ix;
                return true;
            }
        });
        sidoTitles[0][0] += arr[index].DefCnt;
        sidoTitles[1][0] += arr[index].IsolClearCnt;
        sidoTitles[2][0] += arr[index].DeathCnt;
        sidoTitles[0][1] += arr[index + arr.length / 2].DefCnt;
        sidoTitles[1][1] += arr[index + arr.length / 2].IsolClearCnt;
        sidoTitles[2][1] += arr[index + arr.length / 2].DeathCnt;
        content += makeTable(arr, index);
    }
    sidoTitles[0][1] = sidoTitles[0][0] - sidoTitles[0][1];
    sidoTitles[1][1] = sidoTitles[1][0] - sidoTitles[1][1];
    sidoTitles[2][1] = sidoTitles[2][0] - sidoTitles[2][1];
    setTitles(sidoTitles);
    setCure(Math.floor(sidoTitles[1][0] / sidoTitles[0][0] * 10000) / 100);

    if (!(item === "강원도" || item === "제주도" || item === "검역")) {
        for (let i = 0; i < allSido[item].length; i++) {
            sido.append(`<option>${allSido[item][i]}</option>`);
        }
    }
    $('#tables').html(content);
}