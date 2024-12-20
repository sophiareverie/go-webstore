// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"go-store/types"
	"time"
)

func OrderEntry(Products []types.Product) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><link rel=\"stylesheet\" href=\"../assets/styles/styles.css\"><style>\n            .container {\n                display: flex;\n                align-items: flex-start;\n\n            }\n            .form-container {\n                display: flex;\n                flex-direction: column;\n                gap: 5px; \n                max-width: 700px;\n            }\n            label {\n                display: inline-block;\n                margin-bottom: 5px;\n            }\n            input, select, button {\n                max-width: 700px;\n            }\n            main {\n                flex: 1;\n                padding-right: 20px;\n            }\n            aside {\n                flex: 1;\n            }\n            .customer-table {\n                width: 100%;\n                border-collapse: collapse;\n                margin-top: 10px;\n            }\n            .customer-table th, .customer-table td {\n                border: 1px solid #ddd;\n                padding: 8px;\n            }\n            .customer-table tr:hover {\n                background-color: yellow;\n            }\n            .highlight {\n                background-color: #e0f7fa;\n            }\n        </style></head><script src=\"https://unpkg.com/htmx.org\"></script><body><div class=\"container\"><main><form id=\"orderForm\" hx-post=\"/purchasebrief\" hx-trigger=\"submit\" hx-target=\"#asideContent\" hx-swap=\"innerHTML\"><fieldset><legend>Personal Information</legend><div class=\"form-container\"><label for=\"firstName\">First Name*:</label> <input type=\"text\" id=\"firstName\" name=\"firstName\" required pattern=\"[A-Za-z\\s&#39;]+\" onkeyup=\"showHint(this.value, &#39;first&#39;)\"> <label for=\"lastName\">Last Name*:</label> <input type=\"text\" id=\"lastName\" name=\"lastName\" required pattern=\"[A-Za-z\\s&#39;]+\" onkeyup=\"showHint(this.value, &#39;last&#39;)\"> <label for=\"email\">Email*:</label> <input type=\"email\" id=\"email\" name=\"email\" required onkeyup=\"showHint(this.value, &#39;email&#39;)\"></div></fieldset><fieldset><legend>Product Information</legend><div class=\"form-container\"><label for=\"product\">Product*:</label> <select id=\"product\" name=\"product\" required onchange=\"updateAvailableQuantity()\"><option value=\"\">Select an item</option> ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, product := range Products {
			if product.Inactive == 0 {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<option value=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var2 string
				templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(product.Name)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/order_entry.templ`, Line: 86, Col: 64}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var3 string
				templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(product.Name)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/order_entry.templ`, Line: 86, Col: 83}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" - ")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var4 string
				templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%.2f", product.Price))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/order_entry.templ`, Line: 86, Col: 124}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</option>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</select> <label for=\"available\">Available Quantity:</label> <input type=\"text\" id=\"available\" name=\"available\" readonly> <label for=\"quantity\">Quantity*:</label> <input type=\"number\" id=\"quantity\" name=\"quantity\" min=\"1\" required> <input type=\"hidden\" name=\"timestamp\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var5 string
		templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("%d", time.Now().Unix()))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `templates/order_entry.templ`, Line: 97, Col: 110}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"> <button type=\"submit\">Purchase</button></div></fieldset><button type=\"reset\">Clear All Fields</button></form></main><aside id=\"asideContent\"><h3>Customer List</h3><table class=\"customer-table\" id=\"customerTable\"><thead><tr><th>First Name</th><th>Last Name</th><th>Email</th></tr></thead> <tbody id=\"customerList\"></tbody></table><p id=\"noMatchMessage\" style=\"display: none;\">No matching customers found.</p></aside></div><script>\n        document.getElementById('orderForm').addEventListener('htmx:afterRequest', function(event) {\n            if (event.detail.xhr.status === 200) {\n                document.getElementById('orderForm').reset();\n            }\n        });\n\n            function updateAvailableQuantity() {\n                const product = document.getElementById(\"product\").value;\n                const quantityInput = document.getElementById(\"quantity\");\n\n                if (product) {\n                    fetch(`/get_product_quantity?product=${encodeURIComponent(product)}`)\n                        .then(response => response.json())\n                        .then(data => {\n                            const availableQuantity = data.quantity;\n                            document.getElementById(\"available\").value = data.quantity;\n                            quantityInput.max = availableQuantity;\n\n                        })\n                        .catch(error => console.error(\"Error fetching quantity:\", error));\n                } else {\n                    document.getElementById(\"available\").value = \"\";\n                }\n            }\n\n            function showHint(value, fieldType) {\n\n                document.getElementById(\"noMatchMessage\").style.display = \"none\";\n                const tableBody = document.getElementById(\"customerList\");\n    \n\n                if (value.length === 0) {\n                    document.getElementById(\"customerList\").innerHTML = \"\";\n                    document.getElementById(\"noMatchMessage\").style.display = \"none\";\n                    return;\n                }\n\n                fetch(`/get_customers?searchTerm=${encodeURIComponent(value)}`)\n                    .then(response => response.json())\n                    .then(data => {\n\n                        const tableBody = document.getElementById(\"customerList\");\n                        tableBody.innerHTML = \"\";\n                        \n                        if (data.message) {\n                            document.getElementById(\"noMatchMessage\").style.display = \"block\";\n                            return;\n                        }\n\n                        if (data.length > 0) {\n                            document.getElementById(\"noMatchMessage\").style.display = \"none\";\n                            data.forEach(customer => {\n                                const row = document.createElement(\"tr\");\n                                row.innerHTML = `\n                                    <td>${customer.FirstName}</td>\n                                    <td>${customer.LastName}</td>\n                                    <td>${customer.Email}</td>\n                                `;\n                                row.onclick = function() {\n                                    document.getElementById(\"firstName\").value = customer.FirstName;\n                                    document.getElementById(\"lastName\").value = customer.LastName;\n                                    document.getElementById(\"email\").value = customer.Email;\n                                };\n                                tableBody.appendChild(row);\n                            });\n                        } else {\n                            document.getElementById(\"noMatchMessage\").style.display = \"block\";\n                        }\n                    })\n                    .catch(error => console.error(\"Error fetching customers:\", error));\n                }\n\n        </script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
