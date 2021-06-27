const allSido = {
    "경기도": ["서울", "경기", "인천"],
    "강원도": ["강원"],
    "충청도": ["충남", "충북", "세종", "대전"],
    "전라도": ["전북", "전남", "광주"],
    "경상도": ["경북", "대구", "경남", "울산", "부산"],
    "제주도": ["제주"],
    "검역": ["검역"],
}

var lat = 35.8679305;
var lng = 128.7053832;
var map;

let av;
let h;

const getHospital = async () => {
    h = await hItems();

    h.forEach((v, i) => {
        let marker = new kakao.maps.Marker({
            map: map,
            position: new kakao.maps.LatLng(v.lat, v.lng),
            title: v.facilityName,
            image: markerImage,
        });
        markers.push(marker);
        let overlay = new kakao.maps.CustomOverlay({
            position: marker.getPosition(),
        });
        overlaies.push(overlay);
        overlay.setContent(setContent(v, i));
    });

    markers.forEach((v, i) => {
        kakao.maps.event.addListener(v, "click", function () {
            if (selectedOverlay != null) selectedOverlay.setMap(null);
            let or = overlaies[i];
            selectedOverlay = or;
            if (or.getMap() === null) or.setMap(map);
            map.setCenter(
                new kakao.maps.LatLng(
                    v.getPosition().getLat() + 0.05,
                    v.getPosition().getLng()
                )
            );
            map.setLevel(9);
            polyLine.setPath([v.getPosition(), new kakao.maps.LatLng(lat, lng)]);
        });
    });
}

const getAv = async () => {
    av = await avItems();

    setTitles(av);
}

let markerImage = new kakao.maps.MarkerImage(
    "https://cdn.icon-icons.com/icons2/1749/PNG/512/06_113688.png",
    new kakao.maps.Size(32, 32)
);
let userMarkerImage = new kakao.maps.MarkerImage(
    "https://publicdomainvectors.org/photos/noun-project-462",
    new kakao.maps.Size(14, 25)
);

var markers = [];
var overlaies = [];
let selectedOverlay = null;

var container = document.getElementById("map"); //지도를 담을 영역의 DOM 레퍼런스
var map = new kakao.maps.Map(container, {
    center: new kakao.maps.LatLng(lat, lng),
    level: 9,
}); //지도 생성 및 객체 리턴

let polyLine = new kakao.maps.Polyline({
    map: map,
    strokeWeight: 3,
    strokeColor: "red",
    strokeOpacity: 0.8,
    strokeStyle: "solid",
});


if (navigator.geolocation) {
    //위치 정보를 얻기
    navigator.geolocation.getCurrentPosition(function (pos) {
        lat = pos.coords.latitude; // 위도
        lng = pos.coords.longitude; // 경도

        map.setCenter(new kakao.maps.LatLng(lat, lng));

        new kakao.maps.Marker({
            map: map,
            position: new kakao.maps.LatLng(lat, lng),
            image: userMarkerImage,
        });
    }, function () {
        map.setCenter(new kakao.maps.LatLng(lat, lng));

        new kakao.maps.Marker({
            map: map,
            position: new kakao.maps.LatLng(lat, lng),
            image: userMarkerImage,
        });
    });
} else {
    alert("이 브라우저에서는 Geolocation이 지원되지 않습니다.");
    map.setCenter(new kakao.maps.LatLng(lat, lng));
    new kakao.maps.Marker({
        map: map,
        position: new kakao.maps.LatLng(lat, lng),
        image: userMarkerImage,
    });
}

function closeOverlay(overlay) {
    polyLine.setPath(null);
    overlay.setMap(null);
}

function KmOrM(lat1, lng1, lat2, lng2) {
    function deg2rad(deg) {
        return deg * (Math.PI / 180);
    }
    var R = 6371; // Radius of the earth in km
    var dLat = deg2rad(lat2 - lat1); // deg2rad below
    var dLon = deg2rad(lng2 - lng1);
    var a =
        Math.sin(dLat / 2) * Math.sin(dLat / 2) +
        Math.cos(deg2rad(lat1)) *
        Math.cos(deg2rad(lat2)) *
        Math.sin(dLon / 2) *
        Math.sin(dLon / 2);
    var c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
    var d = R * c; // Distance in km
    if (d < 1) d = Math.floor(d * 1000) + "m";
    else d = Math.floor(d * 10) / 10 + "km";
    return d;
}

function distance(lat1, lng1, lat2, lng2) {
    function deg2rad(deg) {
        return deg * (Math.PI / 180);
    }
    var R = 6371; // Radius of the earth in km
    var dLat = deg2rad(lat2 - lat1); // deg2rad below
    var dLon = deg2rad(lng2 - lng1);
    var a =
        Math.sin(dLat / 2) * Math.sin(dLat / 2) +
        Math.cos(deg2rad(lat1)) *
        Math.cos(deg2rad(lat2)) *
        Math.sin(dLon / 2) *
        Math.sin(dLon / 2);
    var c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
    return R * c;
}

function setContent(h, i) {
    return `<div class="overlaybox" onclick="closeOverlay(overlaies[${i}])">
    <div>
    <h4>${h.facilityName}</h4>
    <div>
    <span>${KmOrM(lat, lng, h.lat, h.lng)}</span>
    <span>${h.phoneNumber}</span>
    </div>
    <p>${h.address}</p>
    <div>
    </div>`;
}

function setTitles(titles) {
    for (let i = 0; i < 3; i++) {
        let title = $($('.title')[i]);
        $(title.children('p')[0]).text(makeComma(titles[i]["FirstCnt"]));
        $(title.children('p')[1]).text(makeComma(titles[i]["SecondCnt"]));
    }
}

function makeComma(str) {
    str = String(str);

    return str.replace(/(\d)(?=(?:\d{3})+(?!\d))/g, '$1,');
}

$(function() {
    getAv()
    getHospital()

    $('#page_all').click(async () => {
        await movePage("index")
    });

    $('#page_ageGender').click(async () => {
        await movePage("ageGender")
    });

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
})