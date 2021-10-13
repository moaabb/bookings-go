
let roomID
let CSRFToken

const parseRoomInfo = (id, token) => {
    roomID = id
    CSRFToken = token

}


document.getElementById("check-availability-button").addEventListener("click", function () {
    let html = `
    <form id="check-availability-form" action="" method="" novalidate class="needs-validation">
        <div class="row">
            <div class="col">
                <div class="row" id="reservation-dates-modal">
                    <div class="col">
                        <input autocomplete="off" disabled required class="form-control" type="text" name="start" id="start" placeholder="Arrival">
                    </div>
                    <div class="col">
                        <input autocomplete="off" disabled required class="form-control" type="text" name="end" id="end" placeholder="Departure">
                    </div>
                </div>
            </div>
        </div>
    </form>
    `;
    attention.custom({
            title: 'Choose your dates',
            msg: html,
            icon: "",
            willOpen: () => {
                const elem = document.getElementById("reservation-dates-modal");
                const rp = new DateRangePicker(elem, {
                    format: 'yyyy-mm-dd',
                    showOnFocus: true,
                    minDate: Date.now()
                })
            },
            didOpen: () => {
                document.getElementById("start").removeAttribute("disabled");
                document.getElementById("end").removeAttribute("disabled");
            },
            callback: function(result) {
                console.log("called");

                let form = document.getElementById("check-availability-form");
                let formData = new FormData(form);
                formData.append("csrf_token", CSRFToken);
                formData.append("room_id", roomID);

                fetch('/search-availability-json', {
                    method: "post",
                    body: formData,
                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.ok ) {
                            attention.custom({
                                icon: "success",
                                showConfirmButton: false,
                                msg: '<p>The Room is Available!</p>'
                                + '<p><a href="/book-room?id='
                                + data.room_id
                                + '&s='
                                + data.start_date
                                + '&e='
                                + data.end_date
                                + '" class="btn btn-primary">'
                                + 'Book Now!</a></p>'
                            })
                        } else {
                            attention.error({
                                msg: "No Availability!"
                            })
                        }
                    })
            }
        });
})


