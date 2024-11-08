I chose HTMX and enjoyed the simplicity of using the same kind of posting in my main.go as i did with the other things i implemented. it was as simple as a descriptor in the <form> tag

                <form id="orderForm" hx-post="/purchasebrief" hx-trigger="submit" hx-target="#asideContent" hx-swap="innerHTML">                    <fieldset>


and a bit of scripting. It was nice.