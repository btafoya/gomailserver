#!/usr/bin/env python3
"""
CSS Issues Evaluation for gomailserver UI
Analyzes layout, styling, responsiveness, and visual bugs
"""

from playwright.sync_api import sync_playwright
import json
import os

def evaluate_css_issues(url):
    """Comprehensive CSS issue detection"""

    issues = {
        "layout_issues": [],
        "overflow_issues": [],
        "responsive_issues": [],
        "visual_bugs": [],
        "accessibility_css": [],
        "performance": {},
        "console_errors": []
    }

    screenshots_dir = '/home/btafoya/projects/gomailserver/claudedocs/screenshots'
    os.makedirs(screenshots_dir, exist_ok=True)

    with sync_playwright() as p:
        # Use Firefox since Chromium isn't supported
        browser = p.firefox.launch(headless=True)

        # Test different viewports
        viewports = [
            {"name": "desktop_1920", "width": 1920, "height": 1080},
            {"name": "desktop_1366", "width": 1366, "height": 768},
            {"name": "tablet", "width": 768, "height": 1024},
            {"name": "mobile_large", "width": 414, "height": 896},
            {"name": "mobile_small", "width": 375, "height": 667}
        ]

        for vp in viewports:
            print(f"\n{'='*60}")
            print(f"📱 Testing {vp['name']} ({vp['width']}x{vp['height']})")
            print('='*60)

            context = browser.new_context(
                viewport={'width': vp['width'], 'height': vp['height']},
                user_agent='Mozilla/5.0 (X11; Linux x86_64) Firefox/120.0'
            )
            page = context.new_page()

            # Capture console messages
            console_msgs = []
            page.on('console', lambda msg: console_msgs.append({
                'type': msg.type,
                'text': msg.text
            }))

            try:
                print(f"🌐 Navigating to {url}")
                page.goto(url, wait_until='networkidle', timeout=30000)
                page.wait_for_timeout(2000)  # Wait for Vue to mount

                # Take screenshot
                screenshot_path = f'{screenshots_dir}/css_{vp["name"]}.png'
                page.screenshot(path=screenshot_path, full_page=True)
                print(f"📸 Screenshot saved: {screenshot_path}")

                # === CHECK 1: Horizontal Overflow ===
                overflow_check = page.evaluate('''() => {
                    const body = document.body;
                    const html = document.documentElement;
                    return {
                        hasHorizontalScroll: html.scrollWidth > html.clientWidth,
                        bodyScrollWidth: body.scrollWidth,
                        bodyClientWidth: body.clientWidth,
                        htmlScrollWidth: html.scrollWidth,
                        htmlClientWidth: html.clientWidth,
                        overflowingElements: Array.from(document.querySelectorAll('*')).filter(el => {
                            return el.scrollWidth > el.clientWidth;
                        }).map(el => ({
                            tag: el.tagName,
                            class: el.className,
                            id: el.id,
                            scrollWidth: el.scrollWidth,
                            clientWidth: el.clientWidth
                        })).slice(0, 10)
                    };
                }''')

                if overflow_check['hasHorizontalScroll']:
                    issues['overflow_issues'].append({
                        'viewport': vp['name'],
                        'type': 'horizontal_scroll',
                        'severity': 'high',
                        'details': overflow_check,
                        'description': f"Page has horizontal scrollbar at {vp['width']}px width"
                    })
                    print(f"⚠️  ISSUE: Horizontal overflow detected")
                    print(f"   HTML: {overflow_check['htmlScrollWidth']}px / {overflow_check['htmlClientWidth']}px")

                # === CHECK 2: Layout Issues ===
                layout_check = page.evaluate('''() => {
                    const issues = [];

                    // Check for elements outside viewport
                    document.querySelectorAll('*').forEach(el => {
                        const rect = el.getBoundingClientRect();
                        if (rect.right > window.innerWidth + 50) {
                            issues.push({
                                type: 'element_overflow_right',
                                element: el.tagName + (el.className ? '.' + el.className.split(' ')[0] : ''),
                                right: rect.right,
                                viewportWidth: window.innerWidth
                            });
                        }
                    });

                    // Check for overlapping elements (z-index issues)
                    const elements = Array.from(document.querySelectorAll('*'));
                    const overlaps = [];

                    // Check for fixed/absolute positioning issues
                    elements.forEach(el => {
                        const style = window.getComputedStyle(el);
                        const pos = style.position;
                        if (pos === 'fixed' || pos === 'absolute') {
                            const rect = el.getBoundingClientRect();
                            overlaps.push({
                                element: el.tagName + (el.className ? '.' + el.className.split(' ')[0] : ''),
                                position: pos,
                                zIndex: style.zIndex,
                                top: rect.top,
                                left: rect.left
                            });
                        }
                    });

                    return {
                        overflowingElements: issues.slice(0, 5),
                        positionedElements: overlaps.slice(0, 10)
                    };
                }''')

                if layout_check['overflowingElements']:
                    for elem in layout_check['overflowingElements']:
                        issues['layout_issues'].append({
                            'viewport': vp['name'],
                            'type': 'element_overflow',
                            'severity': 'medium',
                            'element': elem['element'],
                            'details': elem
                        })
                        print(f"⚠️  ISSUE: Element overflow - {elem['element']}")

                # === CHECK 3: Text Rendering Issues ===
                text_issues = page.evaluate('''() => {
                    const problems = [];

                    document.querySelectorAll('*').forEach(el => {
                        const style = window.getComputedStyle(el);
                        const text = el.textContent.trim();

                        if (text.length > 0) {
                            // Check if text is cut off
                            if (style.overflow === 'hidden' && el.scrollHeight > el.clientHeight) {
                                problems.push({
                                    type: 'text_overflow_hidden',
                                    element: el.tagName + (el.className ? '.' + el.className.split(' ')[0] : ''),
                                    text: text.substring(0, 50)
                                });
                            }

                            // Check for very small font sizes
                            const fontSize = parseFloat(style.fontSize);
                            if (fontSize < 12 && el.tagName !== 'SUP' && el.tagName !== 'SUB') {
                                problems.push({
                                    type: 'font_too_small',
                                    element: el.tagName + (el.className ? '.' + el.className.split(' ')[0] : ''),
                                    fontSize: fontSize
                                });
                            }
                        }
                    });

                    return problems.slice(0, 10);
                }''')

                for text_issue in text_issues:
                    issues['visual_bugs'].append({
                        'viewport': vp['name'],
                        'severity': 'medium',
                        **text_issue
                    })

                # === CHECK 4: Missing Styles / Unstyled Elements ===
                style_check = page.evaluate('''() => {
                    const unstyled = [];

                    document.querySelectorAll('button, input, select, textarea').forEach(el => {
                        const style = window.getComputedStyle(el);

                        // Check if element looks default/unstyled
                        const hasCustomStyling =
                            style.backgroundColor !== 'rgba(0, 0, 0, 0)' ||
                            style.border !== 'none' ||
                            style.padding !== '0px' ||
                            el.className.length > 0;

                        if (!hasCustomStyling) {
                            unstyled.push({
                                element: el.tagName,
                                type: el.type || 'N/A',
                                id: el.id || 'N/A',
                                className: el.className || 'N/A'
                            });
                        }
                    });

                    return unstyled.slice(0, 5);
                }''')

                if style_check:
                    for unstyle in style_check:
                        issues['visual_bugs'].append({
                            'viewport': vp['name'],
                            'type': 'possibly_unstyled',
                            'severity': 'low',
                            **unstyle
                        })

                # === CHECK 5: Accessibility CSS Issues ===
                a11y_css = page.evaluate('''() => {
                    const a11yIssues = [];

                    // Check for invisible text
                    document.querySelectorAll('*').forEach(el => {
                        const style = window.getComputedStyle(el);
                        const text = el.textContent.trim();

                        if (text.length > 0) {
                            // Text same color as background
                            if (style.color === style.backgroundColor) {
                                a11yIssues.push({
                                    type: 'invisible_text',
                                    element: el.tagName + (el.className ? '.' + el.className.split(' ')[0] : '')
                                });
                            }

                            // Check opacity
                            if (parseFloat(style.opacity) < 0.5) {
                                a11yIssues.push({
                                    type: 'low_opacity_text',
                                    element: el.tagName + (el.className ? '.' + el.className.split(' ')[0] : ''),
                                    opacity: style.opacity
                                });
                            }
                        }
                    });

                    return a11yIssues.slice(0, 10);
                }''')

                for a11y in a11y_css:
                    issues['accessibility_css'].append({
                        'viewport': vp['name'],
                        'severity': 'high',
                        **a11y
                    })

                # Filter console errors for CSS-related issues
                css_console_errors = [msg for msg in console_msgs if
                    msg['type'] in ['error', 'warning'] and
                    any(keyword in msg['text'].lower() for keyword in ['css', 'style', 'stylesheet', 'font'])]

                if css_console_errors:
                    issues['console_errors'].extend([{
                        'viewport': vp['name'],
                        **msg
                    } for msg in css_console_errors])

                print(f"✅ Viewport {vp['name']} analyzed")

            except Exception as e:
                print(f"❌ Error analyzing {vp['name']}: {e}")
                issues['layout_issues'].append({
                    'viewport': vp['name'],
                    'type': 'evaluation_error',
                    'severity': 'critical',
                    'error': str(e)
                })

            finally:
                context.close()

        browser.close()

    return issues

def generate_report(issues):
    """Generate human-readable CSS issues report"""

    report_lines = [
        "# CSS Issues Evaluation Report",
        f"\n**URL**: http://192.168.25.165:5173/admin/",
        f"\n## Summary\n"
    ]

    # Count issues by severity
    severity_counts = {'critical': 0, 'high': 0, 'medium': 0, 'low': 0}
    for category in issues.values():
        if isinstance(category, list):
            for issue in category:
                sev = issue.get('severity', 'low')
                severity_counts[sev] = severity_counts.get(sev, 0) + 1

    report_lines.append(f"- 🔴 Critical: {severity_counts['critical']}")
    report_lines.append(f"- 🟠 High: {severity_counts['high']}")
    report_lines.append(f"- 🟡 Medium: {severity_counts['medium']}")
    report_lines.append(f"- 🟢 Low: {severity_counts['low']}")

    # Overflow Issues
    if issues['overflow_issues']:
        report_lines.append("\n## 🔄 Overflow Issues\n")
        for issue in issues['overflow_issues']:
            report_lines.append(f"### {issue['viewport'].upper()}")
            report_lines.append(f"**Type**: {issue['type']}")
            report_lines.append(f"**Severity**: {issue['severity']}")
            report_lines.append(f"**Description**: {issue['description']}\n")

    # Layout Issues
    if issues['layout_issues']:
        report_lines.append("\n## 📐 Layout Issues\n")
        for issue in issues['layout_issues']:
            report_lines.append(f"### {issue.get('viewport', 'N/A').upper()}")
            report_lines.append(f"**Element**: `{issue.get('element', 'unknown')}`")
            report_lines.append(f"**Type**: {issue['type']}")
            report_lines.append(f"**Severity**: {issue.get('severity', 'unknown')}\n")

    # Visual Bugs
    if issues['visual_bugs']:
        report_lines.append("\n## 🎨 Visual Bugs\n")
        for issue in issues['visual_bugs'][:10]:  # Limit to 10
            report_lines.append(f"- **{issue['type']}** ({issue.get('viewport', 'N/A')}): {issue.get('element', issue.get('text', 'N/A'))}")

    # Accessibility CSS
    if issues['accessibility_css']:
        report_lines.append("\n## ♿ Accessibility CSS Issues\n")
        for issue in issues['accessibility_css']:
            report_lines.append(f"- **{issue['type']}**: {issue['element']} ({issue['viewport']})")

    # Console Errors
    if issues['console_errors']:
        report_lines.append("\n## 🚨 Console Errors (CSS-related)\n")
        for err in issues['console_errors'][:5]:  # Limit to 5
            report_lines.append(f"- [{err['type']}] {err['text'][:100]}")

    return "\n".join(report_lines)

if __name__ == "__main__":
    print("🚀 Starting CSS Evaluation")
    print("="*60)

    url = "http://192.168.25.165:5173/admin/"

    issues = evaluate_css_issues(url)

    # Save JSON results
    json_path = '/home/btafoya/projects/gomailserver/claudedocs/css_issues.json'
    with open(json_path, 'w') as f:
        json.dump(issues, f, indent=2)

    print(f"\n✅ JSON results saved: {json_path}")

    # Generate markdown report
    report = generate_report(issues)
    report_path = '/home/btafoya/projects/gomailserver/claudedocs/CSS_ISSUES.md'
    with open(report_path, 'w') as f:
        f.write(report)

    print(f"📄 Report saved: {report_path}")
    print("\n" + "="*60)
    print("🎉 CSS Evaluation Complete!")
