export default {
    content: [
        "./internal/views/**/*.html"
    ],
    prefix: "",
    theme: {
        container: {
            center: true,
            padding: "2rem",
        },
        fontFamily: {
            body: [
                "Roboto",
                "system-ui",
                "-apple-system",
                "sans-serif",
            ],
            display: [
                "Raleway",
                "system-ui",
                "-apple-system",
                "sans-serif",
            ],
        },
        extend: {
            colors: {
                background: "hsl(var(--background))",
                foreground: "hsl(var(--foreground))",
                primary: "hsl(var(--primary))",
                secondary: "hsl(var(--secondary))",
                muted: "hsl(var(--muted))",
            },
            screens: {
                xs: "360px",
            },
        },
    },
};
