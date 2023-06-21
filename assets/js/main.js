$(function () {
    /*=======================
                UI Slider Range JS
    =========================*/
    $("#slider-range").slider({
        range: true,
        min: 0,
        max: 2500,
        values: [10, 2500],
        slide: function (event, ui) {
            $("#amount").val("$" + ui.values[0] + " - $" + ui.values[1]);
        }
    });
    $("#amount").val("$" + $("#slider-range").slider("values", 0) +
        " - $" + $("#slider-range").slider("values", 1));

    let domShippingCalculationMsg = $("#shipping-calculation-msg")

    $(".province_id").change(function () {
        provinceID = $(".province_id").val()

        $(".city_id").find("option")
            .remove()
            .end()
            .append('<option value="">Pilih Kota / Kabupaten</option>')

        $.ajax({
            url: "/carts/cities?province_id=" + provinceID,
            method: "GET",
            success: function (result) {
                $.each(result.data, function (i, city) {
                    $(".city_id").append(`<option value="${city.city_id}">${city.city_name}</option>`)
                });
            }
        })
    });

    $(".city_id").change(function () {
        let cityID = $(".city_id").val()
        let courier = $(".courier").val()

        $(".shipping_fee_options").find("option")
            .remove()
            .end()
            .append('<option value="">Pilih Paket</option>')

        $.ajax({
            url: "/carts/calculate-shipping",
            method: "POST",
            data: {
                city_id: cityID,
                courier: courier
            },
            success: function (result) {
                domShippingCalculationMsg.html('');
                $.each(result.data, function (i, shipping_fee_option) {
                    $(".shipping_fee_options").append(`<option value="${shipping_fee_option.service}-${shipping_fee_option.fee}">${shipping_fee_option.fee} (${shipping_fee_option.service})</option>`);
                });
            },
            error: function (e) {
                domShippingCalculationMsg.html(`<div class="alert alert-warning">Perhitungan ongkos kirim gagal!</div>`);
            }
        })
    });

    $(".shipping_fee_options").change(function () {
        let cityID = $(".city_id").val()
        let courier = $(".courier").val()
        let shippingFee = $(this).val();

        $.ajax({
            url: "/carts/apply-shipping",
            method: "POST",
            data: {
                shipping_package: shippingFee.split("-")[0],
                city_id: cityID,
                courier: courier
            },
            success: function (result) {
                $("#grand-total").text(result.data.grand_total)
            },
            error: function (e) {
                domShippingCalculationMsg.html(`<div class="alert alert-warning">Pemilihan paket ongkir gagal!</div>`);
            }
        })
    });
});
