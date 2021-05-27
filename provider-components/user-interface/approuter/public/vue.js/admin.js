/* global Vue axios */ //> from vue.html
const $ = sel => document.querySelector(sel)

const faqs = new Vue({

    el: '#app_edit',

    data: {
        list: [],
        faq: undefined,
        count: { amount: 1, succeeded: '', failed: '' }
    },

    created() {
        let contextUpdateFunction = initialContext => {
            // initially fill list of books
            faqs.fetch("&$orderby=state desc");

            //check if we shall open modal
            var idInModal = LuigiClient.getNodeParams().openInModal;

            if (!!idInModal) {
                LuigiClient.linkManager().openAsModal('/home/edit/' + idInModal, { title: "ID: " + idInModal, size: 'l' });
            }
        }
        LuigiClient.addInitListener(contextUpdateFunction);
        LuigiClient.addContextUpdateListener(contextUpdateFunction);
    },

    methods: {

        search: ({ target: { value: v } }) => faqs.fetch(v && '&$search=' + v),

        async fetch(etc = '') {
            const { data } = await axios.get(`/admin/Faqs?$expand=category,author${etc}`)
            faqs.list = data.value.map(elem => {
                elem.stateRendered = "valid";
                return elem;
            });

        },

        isAnswered(state) {
            return state === 'answered';
        },

        add() {
            LuigiClient.linkManager().openAsModal('/home/add/', { title: "New FAQ", size: 'l' });
        },

        async inspect(eve) {
            const faq = faqs.faq = faqs.list[eve.currentTarget.rowIndex - 1]

            const res = await axios.get(`/admin/Faqs/${faq.ID}?$select=descr,count`);
            Object.assign(faq, res.data);

            LuigiClient.linkManager().openAsModal('/home/edit/' + faq.ID, { title: "ID: " + faq.ID, size: 'l' });
        }
    }
})

