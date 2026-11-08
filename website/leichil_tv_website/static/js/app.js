/*
   PARTICLES
 */

const container = document.getElementById("particles");

if (container) {

    for (let i = 0; i < 35; i++) {

        const particle = document.createElement("span");

        particle.className = "particle";

        particle.style.left = `${Math.random() * 100}%`;
        particle.style.top = `${Math.random() * 100}%`;

        particle.style.animationDelay = `${Math.random() * 8}s`;
        particle.style.animationDuration = `${12 + Math.random() * 10}s`;

        particle.style.opacity = Math.random() * 0.5 + 0.15;

        const size = Math.random() * 4 + 2;

        particle.style.width = `${size}px`;
        particle.style.height = `${size}px`;

        container.appendChild(particle);

    }

}


/*
   MOBILE NAVBAR
 */

const menuToggle = document.getElementById("menuToggle");
const mobileMenu = document.getElementById("mobileMenu");

if (menuToggle && mobileMenu) {

    const closeMenu = () => {

        mobileMenu.classList.remove("active");

        menuToggle.classList.remove("active");

        menuToggle.setAttribute("aria-expanded", "false");

        document.body.style.overflow = "";

    };

    const openMenu = () => {

        mobileMenu.classList.add("active");

        menuToggle.classList.add("active");

        menuToggle.setAttribute("aria-expanded", "true");

        document.body.style.overflow = "hidden";

    };

    menuToggle.addEventListener("click", (e) => {

        e.stopPropagation();

        if (mobileMenu.classList.contains("active")) {

            closeMenu();

        } else {

            openMenu();

        }

    });

    document.addEventListener("click", (e) => {

        if (
            !mobileMenu.contains(e.target) &&
            !menuToggle.contains(e.target)
        ) {

            closeMenu();

        }

    });

    mobileMenu.querySelectorAll("a").forEach(link => {

        link.addEventListener("click", () => {

            closeMenu();

        });

    });

    document.addEventListener("keydown", (e) => {

        if (e.key === "Escape") {

            closeMenu();

        }

    });

    window.addEventListener("resize", () => {

        if (window.innerWidth > 768) {

            closeMenu();

        }

    });

}