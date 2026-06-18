-- +goose Up
-- +goose StatementBegin
CREATE TABLE awareness_campaign_channels
(
    id INT PRIMARY KEY IDENTITY(1,1),

    awareness_campaign_details_id INT NOT NULL,
    awareness_channel_id INT NOT NULL,

    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_awareness_campaign_channels_campaign FOREIGN KEY (awareness_campaign_details_id)
        REFERENCES awareness_campaign_details(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaign_channels_channel FOREIGN KEY (awareness_channel_id)
        REFERENCES awareness_channels(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_awareness_campaign_channels_user_1 FOREIGN KEY (created_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE awareness_campaign_channels;
-- +goose StatementEnd
