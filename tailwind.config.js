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
                "Helvetica Neue",
                "Helvetica",
                "Arial",
                "sans-serif",
            ],
            display: [
                "Raleway",
                "Segoe UI",
                "Tahoma",
                "Geneva",
                "Verdana",
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
