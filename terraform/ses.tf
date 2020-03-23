resource "aws_ses_domain_identity" "sender" {
  domain = var.email_from_domain
}

data "aws_route53_zone" "ses" {
  name         = format("%s.", var.email_from_domain)
  private_zone = false
}

resource "aws_route53_record" "sender_ses_verification_record" {
  zone_id = data.aws_route53_zone.ses.zone_id
  name    = format("_amazonses.%s", aws_ses_domain_identity.sender.id)
  type    = "TXT"
  ttl     = "600"
  records = [aws_ses_domain_identity.sender.verification_token]
}

resource "aws_ses_domain_identity_verification" "sender_ses_verification" {
  domain = aws_ses_domain_identity.sender.id

  depends_on = [aws_route53_record.sender_ses_verification_record]
}

resource "aws_ses_domain_dkim" "sender" {
  domain = aws_ses_domain_identity.sender.id
}

resource "aws_route53_record" "sender_ses_dkim_record" {
  count   = 3
  zone_id = data.aws_route53_zone.ses.zone_id
  name    = format("%s._domainkey.example.com", aws_ses_domain_dkim.sender.dkim_tokens[count.index])
  type    = "CNAME"
  ttl     = "600"
  records = [format("%s.dkim.amazonses.com", aws_ses_domain_dkim.sender.dkim_tokens[count.index])]
}

resource "aws_ses_domain_mail_from" "sender" {
  domain           = aws_ses_domain_identity.sender.id
  mail_from_domain = format("%s.%s", var.email_from, aws_ses_domain_identity.sender.id)

  # Since the SES DNS records can take up to 72 hours I am using `UseDefaultValue`.
  # In most production scenarios using `RejectMessage` will be better to not confuse
  # the end users.
  behavior_on_mx_failure = "UseDefaultValue"
}

resource "aws_route53_record" "sender_ses_domain_mail_from_mx" {
  zone_id = data.aws_route53_zone.ses.zone_id
  name    = aws_ses_domain_identity.sender.id
  type    = "MX"
  ttl     = "600"
  records = [format("10 feedback-smtp.%s.amazonses.com", var.region)]
}

resource "aws_route53_record" "sender_ses_domain_mail_from_txt" {
  zone_id = data.aws_route53_zone.ses.zone_id
  name    = aws_ses_domain_identity.sender.id
  type    = "TXT"
  ttl     = "600"
  records = ["v=spf1 include:amazonses.com -all"]
}

