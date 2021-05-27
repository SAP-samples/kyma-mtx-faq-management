/* global Vue axios */ //> from vue.html
const $ = sel => document.querySelector(sel)

const faqs = new Vue({

    el: '#app_add',

    data: {
        faq: {},
        authors: [],
        categories: [],
        states: [
            'open',
            'answered',
            'duplicate'
        ],
        activeItem: 'home'
    },

    created() {
        let contextUpdateFunction = initialContext => {
            faqs.loadDropdownValues()

        }
        LuigiClient.addInitListener(contextUpdateFunction);
        LuigiClient.addContextUpdateListener(contextUpdateFunction);
    },

    methods: {
        async loadDropdownValues() {
            faqs.faq = {
                title: "",
                descr: "",
                state: "",
                author_ID: "",
                category_ID: "",
            };

            // read the authors
            var res_authors = await axios.get(`/admin/Authors`);
            faqs.authors = res_authors.data.value;
            Object.assign(faqs.authors, res_authors.value);

            // read the authors
            var res_categories = await axios.get(`/admin/Categories`);
            faqs.categories = res_categories.data.value;
            Object.assign(faqs.categories, res_categories.value);
        },

        async save() {
            await axios.post("/admin/Faqs", {
                title: faqs.faq.title,
                descr: faqs.faq.descr,
                answer: faqs.faq.answer,
                state: faqs.faq.state,
                author_ID: faqs.faq.author?.ID,
                category_ID: faqs.faq.category?.ID
            });
            LuigiClient.linkManager().goBack();
        },

        isActive(menuItem) {
            return this.activeItem === menuItem
        },

        setActive(menuItem) {
            this.activeItem = menuItem
        }
    }
})