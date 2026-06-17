# Formatting Review: v1

## Scores

| Dimension | Weight | Score | Weighted |
|---|---|---|---|
| Heading hierarchy | 1x | 9 | 9 |
| Code formatting | 1x | 4 | 4 |
| CTA placement | 2x | 3 | 6 |
| SEO readiness | 1x | 7 | 7 |
| Link strategy | 1x | 3 | 3 |
| Editorial compliance | 2x | 5 | 10 |
| Brand standards | 1x | 6 | 6 |
| Word count | 1x | 8 | 8 |
| **Total** | | | **53 / 100** |

**Normalized score: 5.3 / 10**

## Line-level feedback

- **Line 1**: Title is 65 characters ("Deploying Ory Talos on OpenShift: API key management for AI platforms"). Trim to 50-60 characters for SEO. Suggestion: "Deploying Ory Talos on OpenShift for API key management" (55 chars).
- **Line 3**: First mention of "OpenShift AI" should be "Red Hat OpenShift AI."
- **Line 7**: "TTLs" — expand "TTL" on first use ("time-to-live").
- **Line 9**: "gRPC-Gateway" — expand "gRPC" on first use ("gRPC (Google Remote Procedure Call)") or at minimum note what gRPC-Gateway means for readers unfamiliar with it.
- **Line 19**: Inline backticks around `CGO_ENABLED=0` and `modernc.org/sqlite`. Remove all inline backticks per editorial rules. Use italics or rewrite to avoid monospace formatting in final output.
- **Line 21**: "UBI Dockerfile" — expand "UBI" on first use ("Universal Base Image (UBI)"). Also "Dockerfile" should remain as-is (proper noun).
- **Line 44**: Inline backticks on `ubi-micro`, `chgrp`, `chmod`. Remove backticks.
- **Line 56**: Inline backticks on `no such table: networks`. Remove backticks; rephrase as a quoted error message instead.
- **Line 57**: Inline backticks on `docker-compose.oss.yaml` and `talos migrate up`. Remove backticks.
- **Line 66**: Inline backticks on `talos serve`. Remove backticks.
- **Line 75-78**: Bold + inline content for the list items is fine, but "ConfigMap," "Secret," "PVC," and "Service" should not rely on Kubernetes jargon without context. "PVC" must be expanded on first use ("PersistentVolumeClaim (PVC)").
- **Line 76**: "HMAC" — expand on first use ("Hash-based Message Authentication Code (HMAC)").
- **Line 75**: "CORS" — expand on first use ("Cross-Origin Resource Sharing (CORS)").
- **Line 98**: `GET /health/alive` — inline backticks. Remove.
- **Line 100**: `GET /version` — inline backticks. Remove.
- **Line 102**: `POST /v2alpha1/admin/issuedApiKeys` — inline backticks. Remove.
- **Line 115**: `POST /v2alpha1/admin/apiKeys:verify` — inline backticks. Remove.
- **Line 119**: "sub-millisecond" — consistent with line 7, which is good.
- **Line 123**: Inline backticks on `CGO_ENABLED=0` and `ubi-micro`. Remove.
- **Line 127**: Inline backticks on `/v2alpha1/admin/...`. Remove.
- **Line 129**: Inline backticks on multiple items. Remove.
- **Line 133**: The only link in the entire post goes to github.com. Add at least 2-3 internal links to redhat.com content (e.g., Red Hat OpenShift AI product page, UBI documentation, OpenShift developer resources).
- **Lines 133-151**: CTA appears only at the very end. Add a brief CTA or teaser link near the top of the post (after the intro paragraph) and another mid-article (after the deployment section).

## Editorial compliance checklist

| Rule | Status | Notes |
|---|---|---|
| Sentence case headings | PASS | All headings use sentence case correctly. Capitalized after colon on line 1. |
| Oxford commas | PASS | Consistently used throughout (lines 7, 13, 78). |
| No backticks | FAIL | 20+ instances of inline backticks throughout the post. Every section uses them. |
| Full product name on first mention | FAIL | "OpenShift AI" on line 3 should be "Red Hat OpenShift AI." |
| Lowercase component descriptors | PASS | No issues found. |
| No H1 in body | PASS | All headings are H2. |
| Expand acronyms on first use | FAIL | UBI, TTL, gRPC, CORS, HMAC, PVC, JWT are never expanded. |
| Use contractions | PASS | Good use throughout: "don't," "it's," "you'd," "aren't," "isn't." |
| Numerals in running text | PASS | "1-hour TTL" (line 102), "4 minutes" (line 52), "50 MB" (line 44). |
| No em dashes (or max 1-2) | PASS | No em dashes found. Colons and commas used instead. |

## Summary

The draft has solid structure and clean heading hierarchy. The biggest issues are:

1. **Backticks everywhere (code formatting, score 4):** The post has 20+ inline backtick usages. For Red Hat Developer Blog final output, all inline code references need to be reformatted — use italics, quotation marks for error messages, or rewrite sentences to avoid inline monospace. Code blocks are fine and well-formatted.

2. **No CTAs until the end (CTA placement, score 3):** The "Try it yourself" section at the bottom is the only CTA. Red Hat editorial standards expect a CTA near the top (e.g., link to the OpenShift AI product page or a "get started" link after the intro), a mid-article touchpoint, and the closing CTA. All 3 placements are needed.

3. **No redhat.com links (link strategy, score 3):** The sole external link goes to GitHub. Add internal links to Red Hat product pages, UBI docs, or OpenShift developer resources. This is critical for SEO and editorial compliance.

4. **Unexpanded acronyms (editorial compliance, score 5):** At least 7 acronyms (UBI, TTL, gRPC, CORS, HMAC, PVC, JWT) are used without expansion on first mention. Expand each on first use with the full term followed by the abbreviation in parentheses.

5. **Product name compliance (editorial compliance):** First mention of OpenShift AI must include "Red Hat" prefix.

These are all fixable in v2 without restructuring the content. The underlying writing, heading structure, word count, and use of contractions and Oxford commas are already in good shape.
