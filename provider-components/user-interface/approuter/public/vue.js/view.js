/* global Vue axios */ //> from vue.html
const $ = sel => document.querySelector(sel)

const faqs = new Vue({

    el: '#app_view',

    data: {
        list: [],
        faq: undefined,
        count: { amount: 1, succeeded: '', failed: '' }
    },

    created() {
        let contextUpdateFunction = initialContext => {
            // initially fill list of books
            faqs.fetch()
        }
        LuigiClient.addInitListener(contextUpdateFunction);
        LuigiClient.addContextUpdateListener(contextUpdateFunction);
    },

    methods: {
        search: ({ target: { value: v } }) => faqs.fetch(v && '&$search=' + v),

        async fetch(etc = '') {
            const { data } = await axios.get(`/admin/Faqs?$expand=category,author${etc}`)
            faqs.list = data.value
        },

        async inspect(eve) {
            const faq = faqs.faq = faqs.list[eve.currentTarget.rowIndex - 1]
            const res = await axios.get(`/admin/Faqs/${faq.ID}?$select=descr,count`)
            Object.assign(faq, res.data)
            faqs.count = { count: 1 }
            setTimeout(() => $('form > input').focus(), 111)
        },

        async submitCount() {
            const { faq, count } = faqs, amount = parseInt(faq.count) || 1 // REVISIT: Okra should be less strict
            try {
                const res = await axios.post(`/admin/countFaq`, { amount, faq: faq.ID })
                faqs.count = { amount, succeeded: `Successfully viewed ${amount} times.` }
            } catch (e) {
                faqs.count = { amount, failed: e.response.data.error.message }
            }
        }

    }
})