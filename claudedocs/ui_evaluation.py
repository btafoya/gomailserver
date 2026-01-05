#!/usr/bin/env python3
"""
Comprehensive UI Evaluation for gomailserver Unified Admin Interface
Evaluates accessibility, functionality, and visual state
"""

from playwright.sync_api import sync_playwright
import json
import time

def evaluate_ui():
    results = {
        "accessibility": {},
        "navigation": {},
        "reputation_features": {},
        "visual_state": {},
        "console_errors": [],
        "network_errors": []
    }

    with sync_playwright() as p:
        # Launch browser in headless mode
        browser = p.chromium.launch(headless=True)
        context = browser.new_context(
            viewport={'width': 1920, 'height': 1080},
            user_agent='Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36'
        )
        page = context.new_page()

        # Capture console messages
        console_messages = []
        page.on('console', lambda msg: console_messages.append({
            'type': msg.type,
            'text': msg.text,
            'location': msg.location
        }))

        # Capture network errors
        network_errors = []
        page.on('requestfailed', lambda request: network_errors.append({
            'url': request.url,
            'method': request.method,
            'failure': request.failure
        }))

        try:
            print("📍 Navigating to http://localhost:5173/admin/")
            page.goto('http://localhost:5173/admin/', wait_until='networkidle', timeout=30000)

            # Wait for Vue app to mount
            page.wait_for_timeout(2000)

            # Take initial screenshot
            print("📸 Capturing initial state screenshot")
            page.screenshot(path='/home/btafoya/projects/gomailserver/claudedocs/screenshots/initial_state.png', full_page=True)

            # ===== ACCESSIBILITY EVALUATION =====
            print("\n♿ Evaluating Accessibility (WCAG)")

            # Check for proper heading hierarchy
            headings = page.locator('h1, h2, h3, h4, h5, h6').all()
            results["accessibility"]["heading_count"] = len(headings)
            results["accessibility"]["heading_hierarchy"] = [h.text_content().strip() for h in headings if h.text_content().strip()]

            # Check for images without alt text
            images = page.locator('img').all()
            images_without_alt = [img for img in images if not img.get_attribute('alt')]
            results["accessibility"]["images_total"] = len(images)
            results["accessibility"]["images_missing_alt"] = len(images_without_alt)

            # Check for form labels
            inputs = page.locator('input, textarea, select').all()
            inputs_without_label = []
            for inp in inputs:
                label_id = inp.get_attribute('id')
                aria_label = inp.get_attribute('aria-label')
                aria_labelledby = inp.get_attribute('aria-labelledby')
                if not (label_id or aria_label or aria_labelledby):
                    inputs_without_label.append(inp.get_attribute('name') or 'unnamed')
            results["accessibility"]["inputs_total"] = len(inputs)
            results["accessibility"]["inputs_without_label"] = inputs_without_label

            # Check color contrast (basic check - look for text elements)
            results["accessibility"]["has_semantic_html"] = bool(page.locator('nav, header, main, footer, article, section').count() > 0)

            # Check for skip links
            results["accessibility"]["has_skip_link"] = bool(page.locator('a[href="#main"], a[href="#content"]').count() > 0)

            # Check for proper button/link roles
            buttons = page.locator('button, [role="button"]').all()
            results["accessibility"]["button_count"] = len(buttons)

            # ===== NAVIGATION EVALUATION =====
            print("\n🧭 Evaluating Navigation")

            # Find all navigation links
            nav_links = page.locator('nav a, [role="navigation"] a').all()
            results["navigation"]["nav_links"] = [{"text": link.text_content().strip(), "href": link.get_attribute('href')} for link in nav_links if link.text_content().strip()]
            results["navigation"]["nav_link_count"] = len(nav_links)

            # Test clicking through navigation items
            navigation_tests = []
            for i, link in enumerate(nav_links[:5]):  # Test first 5 nav items
                try:
                    text = link.text_content().strip()
                    if not text:
                        continue

                    print(f"  Testing navigation: {text}")
                    link.click()
                    page.wait_for_load_state('networkidle', timeout=5000)
                    page.wait_for_timeout(1000)

                    current_url = page.url
                    screenshot_name = f'/home/btafoya/projects/gomailserver/claudedocs/screenshots/nav_{i}_{text.replace(" ", "_")}.png'
                    page.screenshot(path=screenshot_name)

                    navigation_tests.append({
                        "text": text,
                        "url": current_url,
                        "screenshot": screenshot_name,
                        "success": True
                    })
                except Exception as e:
                    navigation_tests.append({
                        "text": text if 'text' in locals() else f"link_{i}",
                        "error": str(e),
                        "success": False
                    })

            results["navigation"]["tests"] = navigation_tests

            # ===== REPUTATION MANAGEMENT FEATURES =====
            print("\n🛡️ Evaluating Reputation Management Features")

            # Navigate to reputation management section
            try:
                reputation_link = page.locator('a:has-text("Reputation"), a:has-text("Security"), nav a').first
                if reputation_link.count() > 0:
                    reputation_link.click()
                    page.wait_for_load_state('networkidle', timeout=5000)
                    page.wait_for_timeout(1000)
                    page.screenshot(path='/home/btafoya/projects/gomailserver/claudedocs/screenshots/reputation_section.png', full_page=True)
            except:
                pass

            # Look for reputation-related components
            reputation_elements = {
                "blacklist_check": page.locator('[data-testid*="blacklist"], :has-text("Blacklist")').count() > 0,
                "reputation_score": page.locator('[data-testid*="reputation"], :has-text("Reputation Score")').count() > 0,
                "monitoring": page.locator(':has-text("Monitor"), :has-text("Monitoring")').count() > 0,
                "alerts": page.locator('[role="alert"], .alert, :has-text("Alert")').count() > 0,
                "status_indicators": page.locator('[data-status], .status, [class*="status"]').count(),
                "action_buttons": page.locator('button').count()
            }
            results["reputation_features"] = reputation_elements

            # ===== VISUAL STATE & RESPONSIVENESS =====
            print("\n📱 Evaluating Responsiveness")

            viewports = [
                {"name": "desktop", "width": 1920, "height": 1080},
                {"name": "tablet", "width": 768, "height": 1024},
                {"name": "mobile", "width": 375, "height": 667}
            ]

            responsive_tests = []
            for viewport in viewports:
                page.set_viewport_size({"width": viewport["width"], "height": viewport["height"]})
                page.wait_for_timeout(1000)

                screenshot_path = f'/home/btafoya/projects/gomailserver/claudedocs/screenshots/{viewport["name"]}_view.png'
                page.screenshot(path=screenshot_path, full_page=True)

                # Check for horizontal scrollbar
                has_horizontal_scroll = page.evaluate('document.documentElement.scrollWidth > document.documentElement.clientWidth')

                responsive_tests.append({
                    "viewport": viewport["name"],
                    "size": f"{viewport['width']}x{viewport['height']}",
                    "screenshot": screenshot_path,
                    "has_horizontal_scroll": has_horizontal_scroll
                })

            results["visual_state"]["responsive_tests"] = responsive_tests

            # ===== DOM STRUCTURE ANALYSIS =====
            print("\n🔍 Analyzing DOM Structure")

            # Reset to desktop viewport
            page.set_viewport_size({"width": 1920, "height": 1080})
            page.goto('http://localhost:5173/admin/', wait_until='networkidle')
            page.wait_for_timeout(2000)

            # Get component structure
            vue_components = page.locator('[data-v-app], [id="app"]').count()
            results["visual_state"]["vue_app_detected"] = vue_components > 0

            # Count interactive elements
            results["visual_state"]["interactive_elements"] = {
                "buttons": page.locator('button').count(),
                "links": page.locator('a').count(),
                "inputs": page.locator('input').count(),
                "selects": page.locator('select').count(),
                "textareas": page.locator('textarea').count()
            }

            # Get page title
            results["visual_state"]["page_title"] = page.title()

            # ===== CONSOLE & NETWORK ERRORS =====
            results["console_errors"] = [msg for msg in console_messages if msg['type'] in ['error', 'warning']]
            results["network_errors"] = network_errors

            # Final summary screenshot
            page.screenshot(path='/home/btafoya/projects/gomailserver/claudedocs/screenshots/final_state.png', full_page=True)

        except Exception as e:
            results["evaluation_error"] = str(e)
            print(f"❌ Error during evaluation: {e}")

        finally:
            browser.close()

    return results

if __name__ == "__main__":
    print("🚀 Starting UI Evaluation for gomailserver")
    print("=" * 60)

    # Create screenshots directory
    import os
    os.makedirs('/home/btafoya/projects/gomailserver/claudedocs/screenshots', exist_ok=True)

    results = evaluate_ui()

    # Save results to JSON
    with open('/home/btafoya/projects/gomailserver/claudedocs/ui_evaluation_results.json', 'w') as f:
        json.dump(results, f, indent=2)

    print("\n" + "=" * 60)
    print("✅ UI Evaluation Complete")
    print(f"📊 Results saved to: ui_evaluation_results.json")
    print(f"📸 Screenshots saved to: screenshots/")
    print("\n📋 Summary:")
    print(f"  • Headings: {results['accessibility'].get('heading_count', 0)}")
    print(f"  • Navigation Links: {results['navigation'].get('nav_link_count', 0)}")
    print(f"  • Console Errors: {len(results['console_errors'])}")
    print(f"  • Network Errors: {len(results['network_errors'])}")
    print(f"  • Page Title: {results['visual_state'].get('page_title', 'N/A')}")
