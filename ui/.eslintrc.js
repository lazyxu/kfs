module.exports = {
    extends: [
        'alloy',
        'alloy/react',
    ],
    parserOptions: {
        ecmaVersion: 2022
    },
    env: {
        // Your environments (which contains several predefined global variables)
        //
        browser: true,
        node: true,
        // mocha: true,
        jest: true,
        worker: true
    },
    globals: {
        // Your global variables (setting to false means it's not allowed to be reassigned)
        //
        // myGlobal: false
    },
    rules: {
        // Customize your rules
        "max-params": "off",
        "prefer-arrow-callback": "off",
        "no-unused-vars": "off",
        "react/no-unescaped-entities": "off"
    },
};