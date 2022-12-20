let inputs = $('input[type="button"]')
for (let i = 0, l = inputs.length; i < l; i++) {
    let $input = $(inputs[i])
    for (let $parent = $input.parent(); !$parent.is('body'); $parent = $parent.parent()) {
        if ($parent.is('form')) {
            $input.click(function () {
                let url = $parent.attr("action")
                let method = undefined
                let formData = new FormData()
                $parent.serializeArray().forEach(e => {
                    if (e.name === '_method') {
                        method = e.value
                    } else {
                        formData.append(e.name, e.value);
                    }
                })
                if (!(method)) {
                    method = $parent.attr("method").toUpperCase()
                }

                // application/x-www-form-urlencoded
                $.ajax({
                    url: url,
                    type: method,
                    data: formData,
                    processData: false,
                    contentType: false,
                    success: function () {
                        location.reload()
                    },
                    error: function (e) {
                        // console.error(e)
                    }
                });
            })
            break
        }
    }
}