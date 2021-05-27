/* global Vue axios */ //> from vue.html
const $ = sel => document.querySelector(sel)

const faqs = new Vue({

    el: '#app_edit',

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
            let faqId = LuigiClient.getContext().questionId;
            faqs.loadFaq(faqId);
        }
        LuigiClient.addInitListener(contextUpdateFunction);
        LuigiClient.addContextUpdateListener(contextUpdateFunction);
    },

    methods: {
        async loadFaq(faqId) {
            const res = await axios.get(`/admin/Faqs/${faqId}?$expand=category,author`);

            var data = res.data;
            delete data['@odata.context'];

            faqs.faq = data;

            // read the authors
            var res_authors = await axios.get(`/admin/Authors`);
            faqs.authors = res_authors.data.value;
            Object.assign(faqs.authors, res_authors.value);

            // read the authors
            var res_categories = await axios.get(`/admin/Categories`);
            faqs.categories = res_categories.data.value;
            Object.assign(faqs.categories, res_categories.value);
        },

        async saveQuestion(eve) {
            const res = await axios.patch("/admin/Faqs/" + faqs.faq.ID, {
                ID: faqs.faq.ID,
                descr: faqs.faq.descr
            });
        },
        async saveTitle(eve) {
            const res = await axios.patch("/admin/Faqs/" + faqs.faq.ID, {
                ID: faqs.faq.ID,
                title: faqs.faq.title
            });
        },
        async saveAnswer(eve) {
            const res = await axios.patch("/admin/Faqs/" + faqs.faq.ID, {
                ID: faqs.faq.ID,
                answer: faqs.faq.answer
            });
        },
        async saveQuestion(eve) {
            const res = await axios.patch("/admin/Faqs/" + faqs.faq.ID, {
                ID: faqs.faq.ID,
                descr: faqs.faq.descr
            });
        },
        async saveState(eve) {
            const res = await axios.patch("/admin/Faqs/" + faqs.faq.ID, {
                ID: faqs.faq.ID,
                state: faqs.faq.state
            });
        },

        async deleteQuestion(eve) {
            const res = await axios.delete("/admin/Faqs/" + faqs.faq.ID);
            LuigiClient.linkManager().goBack();
        },

        async saveAuthor(eve) {
            const res = await axios.patch("/admin/Faqs/" + faqs.faq.ID, {
                ID: faqs.faq.ID,
                author_ID: faqs.faq.author.ID
            });
        },

        async saveCategory(eve) {
            const res = await axios.patch("/admin/Faqs/" + faqs.faq.ID, {
                ID: faqs.faq.ID,
                category_ID: faqs.faq.category.ID
            });
        },
        isActive(menuItem) {
            return this.activeItem === menuItem
        },
        setActive(menuItem) {
            this.activeItem = menuItem
        }
    }
})

